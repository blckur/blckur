package oauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/utils"
	"github.com/dropbox/godropbox/errors"
	"github.com/mrjones/oauth"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"math/rand"
	"strconv"
	"time"
)

type Oauth1 struct {
	Type           string
	ConsumerKey    string
	ConsumerSecret string
	ReqTokenUrl    string
	AuthTokenUrl   string
	AccsTokenUrl   string
	CallbackUrl    string
	consumer       *oauth.Consumer
}

func (o *Oauth1) Config() {
	o.consumer = oauth.NewConsumer(
		o.ConsumerKey,
		o.ConsumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   o.ReqTokenUrl,
			AuthorizeTokenUrl: o.AuthTokenUrl,
			AccessTokenUrl:    o.AccsTokenUrl,
		},
	)
}

func (o *Oauth1) Request(db *database.Database, userId bson.ObjectId) (
	url string, err error) {

	coll := db.Tokens()

	reqTokn, url, err := o.consumer.GetRequestTokenAndUrl(o.CallbackUrl)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth.oauth1: Unknown api error"),
		}
		return
	}

	tokn := &Token{
		Id:          reqTokn.Token,
		Type:        o.Type,
		OauthSecret: reqTokn.Secret,
		UserId:      userId,
	}

	err = coll.Insert(tokn)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func (o *Oauth1) Authorize(db *database.Database, token string, code string) (
	client *Oauth1Client, err error) {

	coll := db.Tokens()
	tokn := &Token{}

	err = coll.FindOneId(token, tokn)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	reqTokn := &oauth.RequestToken{
		Token:  tokn.Id,
		Secret: tokn.OauthSecret,
	}

	accessTokn, err := o.consumer.AuthorizeToken(reqTokn, code)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth.oauth1: Unknown api error"),
		}
		return
	}

	acct := &account.Account{
		UserId:    tokn.UserId,
		Type:      o.Type,
		New:       true,
		OauthTokn: accessTokn.Token,
		OauthSec:  accessTokn.Secret,
	}

	client = &Oauth1Client{
		Account: acct,
		Token:   accessTokn.Token,
		Secret:  accessTokn.Secret,
		conf:    o,
	}

	return
}

func (o *Oauth1) NewClient(acct *account.Account) (client *Oauth1Client) {
	client = &Oauth1Client{
		Account: acct,
		Token:   acct.OauthTokn,
		Secret:  acct.OauthSec,
		conf:    o,
	}

	return
}

type Oauth1Client struct {
	Account *account.Account
	Token   string
	Secret  string
	conf    *Oauth1
}

func (c *Oauth1Client) GetJson(url string, userParams map[string]string,
	resp interface{}) (err error) {

	tokn := &oauth.AccessToken{
		Token:  c.Token,
		Secret: c.Secret,
	}

	httpResp, err := c.conf.consumer.Get(url, userParams, tokn)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth.oauth1: Unknown api error"),
		}
		return
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth.oauth1: Unknown parse error"),
		}
		return
	}

	err = json.Unmarshal(body, resp)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth.oauth1: Unknown parse error"),
		}
		return
	}

	return
}

func (c *Oauth1Client) Sign(method string, url string,
	params map[string]string) (sig string) {
	message := method + "&" + utils.Escape(url)
	delimEsc := utils.Escape("&")
	first := true

	addParam := func(key string, val string) {
		if first {
			first = false
			message += "&"
		} else {
			message += delimEsc
		}
		message += utils.Escape(fmt.Sprintf("%s=%s", key, val))
	}

	addParam("oauth_consumer_key", c.conf.ConsumerKey)
	addParam("oauth_nonce", strconv.FormatInt(rand.Int63(), 10))
	addParam("oauth_signature_method", "HMAC-SHA1")
	addParam("oauth_timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	addParam("oauth_token", c.Token)
	addParam("oauth_version", "1.0")

	for key, val := range params {
		addParam(key, val)
	}

	key := utils.Escape(c.conf.ConsumerSecret) + "&" + utils.Escape(c.Secret)

	hashFunc := hmac.New(sha1.New, []byte(key))
	hashFunc.Write([]byte(message))
	rawSignature := hashFunc.Sum(nil)
	sig = base64.StdEncoding.EncodeToString(rawSignature)

	return
}
