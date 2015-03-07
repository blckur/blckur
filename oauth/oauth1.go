package oauth

import (
	"encoding/json"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/errortypes"
	"github.com/dropbox/godropbox/errors"
	"github.com/mrjones/oauth"
	"labix.org/v2/mgo/bson"
	"io/ioutil"
)

type Oauth1 struct {
	Type string
	ConsumerKey string
	ConsumerSecret string
	ReqTokenUrl string
	AuthTokenUrl string
	AccsTokenUrl string
	CallbackUrl string
	consumer *oauth.Consumer
}

func (o *Oauth1) Config() {
	o.consumer = oauth.NewConsumer(
		o.ConsumerKey,
		o.ConsumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl: o.ReqTokenUrl,
			AuthorizeTokenUrl: o.AuthTokenUrl,
			AccessTokenUrl: o.AccsTokenUrl,
		},
	)
}

func (o *Oauth1) Request(db *database.Database, userId bson.ObjectId) (
		url string, err error) {
	coll := db.Tokens()

	reqTokn, url, err := o.consumer.GetRequestTokenAndUrl(o.CallbackUrl)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth: Unknown oauth1 api error"),
		}
		return
	}

	tokn := &Token{
		Id: reqTokn.Token,
		Type: o.Type,
		OauthSecret: reqTokn.Secret,
		UserId: userId,
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
		Token: tokn.Id,
		Secret: tokn.OauthSecret,
	}

	accessTokn, err := o.consumer.AuthorizeToken(reqTokn, code)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth: Unknown oauth1 api error"),
		}
		return
	}

	client = &Oauth1Client{
		UserId: tokn.UserId,
		Token: accessTokn.Token,
		Secret: accessTokn.Secret,
		conf: o,
	}

	return
}

func (o *Oauth1) NewClient(userId bson.ObjectId, token string,
		secret string) (client *Oauth1Client) {
	client = &Oauth1Client{
		UserId: userId,
		Token: token,
		Secret: secret,
		conf: o,
	}

	return
}

type Oauth1Client struct {
	UserId bson.ObjectId
	Token string
	Secret string
	conf *Oauth1
}

func (c *Oauth1Client) GetJson(url string, userParams map[string]string,
		resp interface{}) (err error) {
	tokn := &oauth.AccessToken{
		Token: c.Token,
		Secret: c.Secret,
	}

	httpResp, err := c.conf.consumer.Get(url, userParams, tokn)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth: Unknown oauth1 api error"),
		}
		return
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth: Unknown parse error"),
		}
		return
	}

	err = json.Unmarshal(body, resp)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth: Unknown parse error"),
		}
		return
	}

	return
}
