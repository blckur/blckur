package accounts

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/notification"
	"github.com/blckur/blckur/oauth"
	"github.com/blckur/blckur/settings"
	"github.com/dropbox/godropbox/errors"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strconv"
	"strings"
	"time"
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

type gitHubBackend struct {
	db       *database.Database
	acct     *account.Account
	etag     string
	url      string
	client   *oauth.Oauth2Client
	interval int
	stop     bool
	lastNotf *notification.Notification
}

type gitHubEvent struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Actor struct {
		Login string `json:"login"`
	} `json:"actor"`
	Org struct {
		Login string `json:"login"`
	} `json:"org"`
	CreatedAt string                 `json:"created_at"`
	Payload   map[string]interface{} `json:"payload"`
}

func (g *gitHubBackend) filter(typ string, repo string) bool {
	for _, filter := range g.acct.Filters {
		switch filter.Type {
		case "all", typ:
			return true
		case typ + "_from":
			if strings.Contains(repo, filter.Value) {
				return true
			}
		}
	}

	return false
}

func (g *gitHubBackend) parse(evt *gitHubEvent) (err error) {
	timestamp, err := time.Parse("2006-01-02T15:04:05Z", evt.CreatedAt)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("account.github: Failed to parse timestamp")
		return
	}

	if g.lastNotf == nil || evt.Id == g.lastNotf.RemoteId ||
		timestamp.Before(g.lastNotf.Timestamp) {

		notf := &notification.Notification{
			UserId:    g.acct.UserId,
			AccountId: g.acct.Id,
			RemoteId:  evt.Id,
			Timestamp: timestamp,
		}

		err = notf.Initialize(g.db)
		if err != nil {
			return
		}
	}

	switch evt.Type {
	case "IssuesEvent", "IssueCommentEvent":
		action := evt.Payload["action"].(string)
		issue := evt.Payload["issue"].(map[string]interface{})
		user := issue["user"].(map[string]interface{})
		from := user["login"].(string)
		title := issue["title"].(string)
		link := issue["url"].(string)
		repo := evt.Payload["repository"].(map[string]interface{})
		repoName := repo["full_name"].(string)

		var typ string
		var subject string

		if evt.Type == "IssueCommentEvent" {
			if action != "created" {
				return
			}

			typ = "issue_comment"
			subject = fmt.Sprintf("New issue comment in %s", repoName)
		} else {
			switch action {
			case "assigned", "unassigned":
				if evt.Payload["assignee"] != g.acct.Identity {
					return
				}
				typ = "issue_" + action
			case "opened":
				typ = "issue_opened"
			case "closed":
				typ = "issue_closed"
			case "reopened":
				typ = "issue_reopened"
			default:
				return
			}

			subject = fmt.Sprintf("Issue %s in %s", action, repoName)
		}

		if !g.filter(typ, repo) {
			return
		}

		notf := &notification.Notification{
			UserId:    g.acct.UserId,
			AccountId: g.acct.Id,
			RemoteId:  evt.Id,
			Timestamp: timestamp,
			Type:      typ,
			Resource:  repoName,
			Origin:    from,
			Link:      link,
			Subject:   subject,
			Body:      title,
		}

		_ = notf
	}

	return
}

func (g *gitHubBackend) sync() {
	req, err := http.NewRequest("GET", g.url, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("account.github: Stream request init error")
		return
	}

	if g.etag != "" {
		req.Header.Add("If-None-Match", g.etag)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("account.github: Stream request error")
		return
	}
	defer resp.Body.Close()

	intervalStr := resp.Header.Get("X-Poll-Interval")
	if intervalStr != "" {
		g.interval, err = strconv.Atoi(intervalStr)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("account.github: Failed to parse poll interval")
			g.interval = 60
		}
	} else {
		g.interval = 60
	}

	if resp.StatusCode == 304 {
		return
	}

	g.etag = resp.Header.Get("ETag")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "accounts.github: Unknown parse error"),
		}
		return
	}

	data := []*gitHubEvent{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "accounts.github: Unknown parse error"),
		}
		return
	}

	for _, evt := range data {
		err = g.parse(evt)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("account.github: Failed to parse event")
		}
	}
}

func (g *gitHubBackend) Run() {
	g.url = fmt.Sprintf("https://api.github.com/users/%s/received_events",
		g.acct.Identity)
	g.client = gitHubConf.NewClient(g.acct)
	g.interval = 0

	for {
		if g.interval > 0 {
			time.Sleep(time.Duration(g.interval) * time.Second)
		}

		if g.stop {
			return
		}

		g.sync()
	}
}

func (g *gitHubBackend) Stop() {
	g.stop = true
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
