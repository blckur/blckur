package account

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/oauth"
	"github.com/ChimeraCoder/anaconda"
	"labix.org/v2/mgo/bson"
)

var (
	twitterConf *oauth.Oauth1
)

type Twitter struct {
	Account `bson:",inline"`
}

func (t *Twitter) NewClient() (client *oauth.Oauth1Client) {
	client = twitterConf.NewClient(t.UserId, t.OauthTokn, t.OauthSec)
	return
}

func (t *Twitter) NewApiClient() (client *anaconda.TwitterApi) {
	client = anaconda.NewTwitterApi(t.OauthTokn, t.OauthSec)
	return
}

func (t *Twitter) Update() (err error) {
	client := t.NewClient()

	data := struct {
		IdStr string `json:"id_str"`
		ScreenName string `json:"screen_name"`
	}{}

	err = client.GetJson(
		"https://api.twitter.com/1.1/account/verify_credentials.json",
		nil, &data)
	if err != nil {
		return
	}

	t.Identity = "@" + data.ScreenName
	t.IdentityId = data.IdStr

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
		acct *Twitter, err error) {
	coll := db.Accounts()

	client, err := twitterConf.Authorize(db, token, code)
	if err != nil {
		return
	}

	acct = &Twitter{
		Account{
			UserId: client.UserId,
			Type: "twitter",
			OauthTokn: client.Token,
			OauthSec: client.Secret,
			coll: coll,
		},
	}

	err = acct.Update()
	if err != nil {
		return
	}

	err = coll.Insert(acct)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func updateTwitter() {
	anaconda.SetConsumerKey(settings.Twitter.ConsumerKey)
	anaconda.SetConsumerSecret(settings.Twitter.ConsumerSecret)

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
