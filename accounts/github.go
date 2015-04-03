package accounts

import (
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/oauth"
	"github.com/blckur/blckur/settings"
	"labix.org/v2/mgo/bson"
)

var (
	gitHubConf *oauth.Oauth2
)

func init() {
	account.Register("github", "GitHub", OAUTH2,
		GitHubAuth{}, GitHubClient{},
		[]*account.FilterType{
			&account.FilterType{
				Label: "All new events",
				Type:  "all",
			},
		}, func() {
			messenger.Register("settings", "github",
				func(_ *messenger.Message) {
					updateGitHub()
				})
			updateGitHub()
		})
}

type GitHubClient struct {
	acct *account.Account
}

func (g *GitHubClient) SetAccount(acct *account.Account) {
	g.acct = acct
}

func (g *GitHubClient) Update(db *database.Database) (err error) {
	client := gitHubConf.NewClient(g.acct)

	data := struct {
		Login string `json:"login"`
	}{}

	err = client.GetJson("https://api.github.com/user", &data)
	if err != nil {
		return
	}

	g.acct.Identity = data.Login

	return
}

func (g *GitHubClient) Sync(db *database.Database) (err error) {
	return
}

type GitHubAuth struct{}

func (g *GitHubAuth) Request(db *database.Database,
	userId bson.ObjectId) (url string, err error) {

	url, err = gitHubConf.Request(db, userId)
	if err != nil {
		return
	}

	return
}

func (g *GitHubAuth) Authorize(db *database.Database, state string,
	code string) (acct *account.Account, err error) {

	coll := db.Accounts()

	auth, err := gitHubConf.Authorize(db, state, code)
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

func updateGitHub() {
	gitHubConf = &oauth.Oauth2{
		Type:         "github",
		ClientId:     settings.GitHub.ClientId,
		ClientSecret: settings.GitHub.ClientSecret,
		CallbackUrl:  settings.System.Domain + "/callback/github",
		AuthUrl:      "https://github.com/login/oauth/authorize",
		TokenUrl:     "https://github.com/login/oauth/access_token",
		Scopes: []string{
			"repo:status",
			"notifications",
		},
	}
	gitHubConf.Config()
}
