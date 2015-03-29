package accounts

import (
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/oauth"
	"labix.org/v2/mgo/bson"
)

var (
	digitalOceanConf *oauth.Oauth2
)

func init() {
	account.Register("digitalocean", "DigitalOcean", OAUTH2,
	DigitalOceanAuth{}, DigitalOceanClient{},
	[]*account.FilterType{
		&account.FilterType{
			Label: "All new events",
			Type: "all",
		},
		&account.FilterType{
			Label: "Droplet created",
			Type: "create",
		},
		&account.FilterType{
			Label: "Droplet destroyed",
			Type: "destroy",
		},
		&account.FilterType{
			Label: "Droplet password reset",
			Type: "password_reset",
		},
	}, func() {
		messenger.Register("settings", "digitalocean",
				func(_ *messenger.Message) {
			updateDigitalOcean()
		})
		updateDigitalOcean()
	})
}

type DigitalOceanClient struct {
	acct *account.Account
}

func (d *DigitalOceanClient) setAccount(acct *account.Account) {
	d.acct = acct
}

func (d *DigitalOceanClient) Update(db *database.Database) (err error) {
	client := digitalOceanConf.NewClient(d.acct)
	err = client.Refresh(db)
	if err != nil {
		return
	}

	data := struct {
		Account struct {
			Email string `json:"email"`
		} `json:"account"`
	}{}

	err = client.GetJson(
		"https://api.digitalocean.com/v2/account", &data)
	if err != nil {
		return
	}

	d.acct.Identity = data.Account.Email

	return
}

func (d *DigitalOceanClient) Sync(db *database.Database) (err error) {
	return
}

type DigitalOceanAuth struct {
}

func (d *DigitalOceanAuth) Request(db *database.Database,
		userId bson.ObjectId) (url string, err error) {
	url, err = digitalOceanConf.Request(db, userId)
	if err != nil {
		return
	}

	return
}

func (d *DigitalOceanAuth) Authorize(db *database.Database, state string,
		code string) (acct *account.Account, err error) {
	coll := db.Accounts()

	auth, err := digitalOceanConf.Authorize(db, state, code)
	if err != nil {
		return
	}

	client, err := auth.Account.GetClient()
	if err != nil {
		return
	}

	err = client.Update(db)
	if err != nil {
		return
	}

	err = coll.Insert(auth.Account)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func updateDigitalOcean() {
	digitalOceanConf = &oauth.Oauth2{
		Type: "digitalocean",
		ClientId: settings.DigitalOcean.ClientId,
		ClientSecret: settings.DigitalOcean.ClientSecret,
		CallbackUrl: settings.System.Domain + "/callback/digitalocean",
		AuthUrl: "https://cloud.digitalocean.com/v1/oauth/authorize",
		TokenUrl: "https://cloud.digitalocean.com/v1/oauth/token",
		Scopes: []string{
			"read",
		},
	}
	digitalOceanConf.Config()
}
