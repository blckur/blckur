package oauth

import (
	"encoding/json"
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/utils"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"golang.org/x/oauth2"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"net/http"
)

type Oauth2 struct {
	Type         string
	ClientId     string
	ClientSecret string
	CallbackUrl  string
	AuthUrl      string
	TokenUrl     string
	Scopes       []string
	conf         *oauth2.Config
}

func (o *Oauth2) Config() {
	o.conf = &oauth2.Config{
		ClientID:     o.ClientId,
		ClientSecret: o.ClientSecret,
		RedirectURL:  o.CallbackUrl,
		Scopes:       o.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  o.AuthUrl,
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
			errors.Wrap(err, "oauth.oauth2: Unknown api error"),
		}
		return
	}

	tokn := &Token{
		Id:     state,
		Type:   o.Type,
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
			errors.Wrap(err, "oauth.oauth2: Unknown api error"),
		}
		return
	}

	acct := &account.Account{
		UserId:        tokn.UserId,
		Type:          o.Type,
		New:           true,
		Oauth2AccTokn: accessTokn.AccessToken,
		Oauth2RefTokn: accessTokn.RefreshToken,
		Oauth2Exp:     accessTokn.Expiry,
	}

	client = &Oauth2Client{
		Account: acct,
		Token:   *accessTokn,
		client:  o.conf.Client(oauth2.NoContext, accessTokn),
		conf:    o,
	}

	return
}

func (o *Oauth2) NewClient(acct *account.Account) (client *Oauth2Client) {
	tokn := &oauth2.Token{
		AccessToken:  acct.Oauth2AccTokn,
		TokenType:    "Bearer",
		RefreshToken: acct.Oauth2RefTokn,
		Expiry:       acct.Oauth2Exp,
	}

	client = &Oauth2Client{
		Account: acct,
		client:  o.conf.Client(oauth2.NoContext, tokn),
		conf:    o,
	}

	return
}

type Oauth2Client struct {
	oauth2.Token
	client  *http.Client
	conf    *Oauth2
	Account *account.Account
}

func (c *Oauth2Client) Refresh(db *database.Database) (err error) {
	refreshed, err := c.Check()
	if err != nil {
		return
	}

	if !refreshed {
		return
	}

	coll := db.Accounts()

	c.Account.Oauth2AccTokn = c.AccessToken
	c.Account.Oauth2RefTokn = c.RefreshToken
	c.Account.Oauth2Exp = c.Expiry

	fields := set.NewSet("oauth2_acc_tokn", "oauth2_ref_tokn", "oauth2_exp")

	err = coll.CommitFields(c.Account.Id, c.Account, fields)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func (c *Oauth2Client) Check() (refreshed bool, err error) {
	tokn, err := c.client.Transport.(*oauth2.Transport).Source.Token()
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth.oauth2: Unknown token error"),
		}
		return
	}

	refreshed = tokn.AccessToken != c.AccessToken
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
			errors.Wrap(err, "oauth.oauth2: Unknown api error"),
		}
		return
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth.oauth2: Unknown parse error"),
		}
		return
	}

	err = json.Unmarshal(body, resp)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "oauth.oauth2: Unknown parse error"),
		}
		return
	}

	return
}
