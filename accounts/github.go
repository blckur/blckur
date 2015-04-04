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
			&account.FilterType{
				Label: "All issues opened",
				Type:  "issue_opened",
			},
			&account.FilterType{
				Label:       "Issues opened matching repository",
				Type:        "issue_opened_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial repository name",
				ValueHolder: "Repository name",
			},
			&account.FilterType{
				Label: "All issues closed",
				Type:  "issue_closed",
			},
			&account.FilterType{
				Label:       "Issues closed matching repository",
				Type:        "issue_closed_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial repository name",
				ValueHolder: "Repository name",
			},
			&account.FilterType{
				Label: "All issue comments",
				Type:  "issue_comment",
			},
			&account.FilterType{
				Label:       "Issue comments matching repository",
				Type:        "issue_comment_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial repository name",
				ValueHolder: "Repository name",
			},
			&account.FilterType{
				Label: "All issues assigned to you",
				Type:  "issue_assigned",
			},
			&account.FilterType{
				Label:       "Issues assigned matching repository",
				Type:        "issue_assigned_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial repository name",
				ValueHolder: "Repository name",
			},
			&account.FilterType{
				Label: "All issues unassigned to you",
				Type:  "issue_unassigned",
			},
			&account.FilterType{
				Label:       "Issues unassigned matching repository",
				Type:        "issue_unassigned_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial repository name",
				ValueHolder: "Repository name",
			},
			&account.FilterType{
				Label: "All issues closed",
				Type:  "issue_closed",
			},
			&account.FilterType{
				Label:       "All issues closed matching repository",
				Type:        "issue_closed_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial repository name",
				ValueHolder: "Repository name",
			},
			&account.FilterType{
				Label: "All issues reopened",
				Type:  "issue_reopened",
			},
			&account.FilterType{
				Label:       "All issues reopened matching repository",
				Type:        "issue_reopened_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial repository name",
				ValueHolder: "Repository name",
			},
			&account.FilterType{
				Label: "",
				Type:  "",
			},
			&account.FilterType{
				Label:       "",
				Type:        "",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial repository name",
				ValueHolder: "Repository name",
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
