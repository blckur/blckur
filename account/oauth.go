package account

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/oauth"
	"github.com/dropbox/godropbox/container/set"
)

func oauth2Refresh(db *database.Database, acct *Account,
		client *oauth.Oauth2Client) (err error) {
	refreshed, err := client.Check()
	if err != nil {
		return
	}

	if !refreshed {
		return
	}

	coll := db.Accounts()

	acct.Oauth2AccTokn = client.AccessToken
	acct.Oauth2RefTokn = client.RefreshToken
	acct.Oauth2Exp = client.Expiry

	fields := set.NewSet("oauth2_acc_tokn", "oauth2_ref_tokn", "oauth2_exp")

	err = coll.CommitFields(acct.Id, acct, fields)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}
