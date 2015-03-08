package oauth

import (
	"encoding/json"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/utils"
	"github.com/dropbox/godropbox/errors"
	"golang.org/x/oauth2"
	"labix.org/v2/mgo/bson"
	"io/ioutil"
	"time"
	"net/http"
)

type Oauth2 struct {
	Type string
	ClientId string
	ClientSecret string
	CallbackUrl string
	AuthUrl string
	TokenUrl string
	Scopes []string
	conf *oauth2.Config
}

func (o *Oauth2) Config() {
	o.conf = &oauth2.Config{
		ClientID: o.ClientId,
		ClientSecret: o.ClientSecret,
		RedirectURL: o.CallbackUrl,
		Scopes: o.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL: o.AuthUrl,
			TokenURL: o.TokenUrl,
		},
	}
}

func (o *Oauth2) Request(db *database.Database, userId bson.ObjectId) (
		url string, err error) {
	coll := db.Tokens()
	state := utils.RandStr(32)

	url = o.conf.AuthCodeURL(state, oauth2.AccessTypeOffline,
		oauth2.ApprovalForce)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "account: Unknown oauth2 api error"),
		}
		return
	}

	tokn := &Token{
		Id: state,
		Type: o.Type,
		UserId: userId,
	}

	err = coll.Insert(tokn)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func (o *Oauth2) Authorize(db *database.Database, state string, code string) (
		client *Oauth2Client, err error) {
	coll := db.Tokens()
	tokn := &Token{}

	err = coll.FindOneId(state, tokn)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	accessTokn, err := o.conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "account: Unknown google api error"),
		}
		return
	}

	client = &Oauth2Client{
		Token: accessTokn,
		UserId: tokn.UserId,
		client: o.conf.Client(oauth2.NoContext, accessTokn),
		conf: o,
	}

	return
}

func (o *Oauth2) NewClient(userId bson.ObjectId, accessToken string,
		refreshToken string, expiry time.Time) (client *Oauth2Client) {
	tokn := &oauth2.Token{
		AccessToken: accessToken,
		TokenType: "Bearer",
		RefreshToken: refreshToken,
		Expiry: expiry,
	}

	client = &Oauth2Client{
		Token: tokn,
		UserId: userId,
		client: o.conf.Client(oauth2.NoContext, tokn),
		conf: o,
	}

	return
}

type Oauth2Client struct {
	*oauth2.Token
	UserId bson.ObjectId
	client *http.Client
	conf *Oauth2
}

func (c *Oauth2Client) Check() (refreshed bool, err error) {
	tokn, err := c.client.Transport.(*oauth2.Transport).Source.Token()
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth: Unknown oauth2 token error"),
		}
		return
	}

	refreshed = !(tokn.AccessToken != c.AccessToken)
	if refreshed {
		c.AccessToken = tokn.AccessToken
		c.RefreshToken = tokn.RefreshToken
		c.Expiry = tokn.Expiry
	}

	return
}

func (c *Oauth2Client) GetJson(url string, resp interface{}) (err error) {
	httpResp, err := c.client.Get(url)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth: Unknown oauth2 api error"),
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
