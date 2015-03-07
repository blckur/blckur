package account

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/oauth"
	"labix.org/v2/mgo/bson"
)

var (
	twitterConf *oauth.Oauth1
)

type Twitter struct {
	*Account
}

func (t *Twitter) Update() (err error) {
	client := twitterConf.NewClient(t.UserId, t.OauthToken, t.OauthSecret)

	data := &struct {
		ScreenName string `json:"screen_name"`
	}{}

	err = client.GetJson(
		"https://api.twitter.com/1.1/account/settings.json", nil, data)
	if err != nil {
		return
	}

	t.Identity = "@" + data.ScreenName

	return
}

func ReqTwitter(db *database.Database, userId bson.ObjectId) (
		url string, err error) {
	url, err = twitterConf.Request(db, userId)
	if err != nil {
		return
	}

	return
}

func AuthTwitter(db *database.Database, token string, code string) (
		acct *Account, err error) {
	coll := db.Accounts()

	client, err := twitterConf.Authorize(db, token, code)
	if err != nil {
		return
	}

	data := &struct {
		ScreenName string `json:"screen_name"`
	}{}

	err = client.GetJson(
		"https://api.twitter.com/1.1/account/settings.json", nil, data)
	if err != nil {
		return
	}

	acct = &Account{
		UserId: client.UserId,
		Type: "twitter",
		Identity: "@" + data.ScreenName,
		OauthToken: client.Token,
		OauthSecret: client.Secret,
		coll: coll,
	}

	err = coll.Insert(acct)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func updateTwitter() {
	conf := &oauth.Oauth1{
		Type: "twitter",
		ConsumerKey: settings.Twitter.ConsumerKey,
		ConsumerSecret: settings.Twitter.ConsumerSecret,
		ReqTokenUrl: "https://api.twitter.com/oauth/request_token",
		AuthTokenUrl: "https://api.twitter.com/oauth/authorize",
		AccsTokenUrl: "https://api.twitter.com/oauth/access_token",
		CallbackUrl: settings.System.Domain + "/callback/twitter",
	}
	conf.Config()

	twitterConf = conf
}

func InitTwitter() {
	requires.After("settings")
	requires.Before("messenger")

	messenger.Register("settings", "twitter", func(_ *messenger.Message) {
		updateTwitter()
	})
	updateTwitter()

	requires.Register("account")
}
