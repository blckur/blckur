package account

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/utils"
	"github.com/blckur/blckur/messenger"
	"github.com/dropbox/godropbox/errors"
	"golang.org/x/oauth2"
	"labix.org/v2/mgo/bson"
)

var (
	gmailConf *oauth2.Config
)

type Gmail struct {
	*Account
}

func AuthGmail(db *database.Database, userId bson.ObjectId) (
		url string, err error) {
	coll := db.Tokens()
	state := utils.RandStr(32)

	url = gmailConf.AuthCodeURL(state, oauth2.AccessTypeOnline)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "account: Unknown twitter api error"),
		}
		return
	}

	tokn := &Token{
		Token: state,
		UserId: userId,
	}

	err = coll.Insert(tokn)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func updateGmail() {
	gmailConf = &oauth2.Config{
		ClientID: "",
		ClientSecret: "",
		RedirectURL: "",
		Scopes: []string{
			"https://www.googleapis.com/auth/gmail.readonly",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL: "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
}

func InitGmail() {
	requires.After("settings")
	requires.Before("messenger")

	messenger.Register("settings", "gmail", func(_ *messenger.Message) {
		updateGmail()
	})
	updateGmail()

	requires.Register("account")
}
