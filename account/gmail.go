package account

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/oauth"
	"labix.org/v2/mgo/bson"
)

var (
	gmailConf *oauth.Oauth2
)

type Gmail struct {
	*Account
}

func ReqGmail(db *database.Database, userId bson.ObjectId) (
		url string, err error) {
	url, err = gmailConf.Request(db, userId)
	if err != nil {
		return
	}

	return
}

func AuthGmail(db *database.Database, state string, code string) (
		acct *Account, err error) {
	coll := db.Accounts()

	client, err := gmailConf.Authorize(db, state, code)
	if err != nil {
		return
	}

	data := &struct {
		EmailAddress string `bson:"emailAddress"`
	}{}

	err = client.GetJson(
		"https://www.googleapis.com/gmail/v1/users/me/profile", data)
	if err != nil {
		return
	}

	acct = &Account{
		UserId: client.UserId,
		Type: "gmail",
		Identity: data.EmailAddress,
		Oauth2AccessToken: client.AccessToken,
		Oauth2RefreshToken: client.RefreshToken,
		Oauth2Expiry: client.Expiry,
		coll: coll,
	}

	err = coll.Insert(acct)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func updateGmail() {
	gmailConf = &oauth.Oauth2{
		Type: "gmail",
		ClientId: settings.Google.ClientId,
		ClientSecret: settings.Google.ClientSecret,
		CallbackUrl: settings.System.Domain + "/callback/gmail",
		AuthUrl: "https://accounts.google.com/o/oauth2/auth",
		TokenUrl: "https://www.googleapis.com/oauth2/v3/token",
		Scopes: []string{
			"https://www.googleapis.com/auth/gmail.readonly",
		},
	}
	gmailConf.Config()
}

func InitGmail() {
	requires.After("settings")
	requires.Before("messenger")

	messenger.Register("settings", "google", func(_ *messenger.Message) {
		updateGmail()
	})
	updateGmail()

	requires.Register("account")
}
