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
	"github.com/blckur/blckur/stream"
	"github.com/blckur/blckur/utils"
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
				Label: "Repository forked",
				Type:  "fork",
			},
			&account.FilterType{
				Label:       "Repository forked matching repository",
				Type:        "fork_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial repository name",
				ValueHolder: "Repository name",
			},
			&account.FilterType{
				Label: "Comment on commit",
				Type:  "commit_comment",
			},
			&account.FilterType{
				Label:       "Comment on commit matching repository",
				Type:        "commit_comment_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial repository name",
				ValueHolder: "Repository name",
			},
			&account.FilterType{
				Label: "All pull requests opened",
				Type:  "pull_opened",
			},
			&account.FilterType{
				Label:       "Pull request opened matching repository",
				Type:        "pull_opened_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial repository name",
				ValueHolder: "Repository name",
			},
			&account.FilterType{
				Label: "All pull requests closed",
				Type:  "pull_closed",
			},
			&account.FilterType{
				Label:       "Pull requests closed matching repository",
				Type:        "pull_closed_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial repository name",
				ValueHolder: "Repository name",
			},
			&account.FilterType{
				Label: "All pull requests reopened",
				Type:  "pull_reopened",
			},
			&account.FilterType{
				Label:       "Pull requests reopened matching repository",
				Type:        "pull_reopened_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial repository name",
				ValueHolder: "Repository name",
			},
			&account.FilterType{
				Label: "All pull request comments",
				Type:  "pull_comment",
			},
			&account.FilterType{
				Label:       "Pull request comments matching repository",
				Type:        "pull_comment_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial repository name",
				ValueHolder: "Repository name",
			},
			&account.FilterType{
				Label: "All pull requests assigned to you",
				Type:  "pull_assigned",
			},
			&account.FilterType{
				Label:       "Pull request assigned matching repository",
				Type:        "pull_assigned_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial repository name",
				ValueHolder: "Repository name",
			},
			&account.FilterType{
				Label: "All pull requests unassigned to you",
				Type:  "issue_unassigned",
			},
			&account.FilterType{
				Label:       "Pull requests unassigned matching repository",
				Type:        "pull_unassigned_from",
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
	backend := &gitHubBackend{
		db:   db,
		acct: g.acct,
	}
	stream := stream.NewStream(db, g.acct.Id, backend)

	err = stream.Start()
	if err != nil {
		return
	}

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

func (g *gitHubBackend) parse(evt *gitHubEvent, force bool) (
	notf *notification.Notification, stop bool, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = &ParseError{
				errors.New(fmt.Sprint(r)),
			}
		}
	}()

	stop = false

	timestamp, err := time.Parse("2006-01-02T15:04:05Z", evt.CreatedAt)
	if err != nil {
		return
	}

	if force {
		notf = &notification.Notification{
			UserId:    g.acct.UserId,
			AccountId: g.acct.Id,
			RemoteId:  evt.Id,
			Timestamp: timestamp,
		}
	}

	if g.lastNotf == nil || evt.Id == g.lastNotf.RemoteId ||
		timestamp.Before(g.lastNotf.Timestamp) {

		stop = true
		return
	}

	switch evt.Type {
	case "IssuesEvent", "IssueCommentEvent":
		action := evt.Payload["action"].(string)
		issue := evt.Payload["issue"].(map[string]interface{})
		user := issue["user"].(map[string]interface{})
		from := user["login"].(string)
		title := issue["title"].(string)
		link := issue["html_url"].(string)
		repo := evt.Repo.Name

		if len(title) > 140 {
			title = title[:140]
		}

		var typ string
		var subject string

		if evt.Type == "IssueCommentEvent" {
			if action != "created" {
				return
			}

			typ = "issue_comment"
			subject = fmt.Sprintf("New issue comment in %s", repo)
		} else {
			switch action {
			case "assigned", "unassigned":
				if issue["assignee"].(string) != g.acct.Identity {
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

			subject = fmt.Sprintf("Issue %s in %s", action, repo)
		}

		if !g.filter(typ, repo) {
			return
		}

		notf = &notification.Notification{
			UserId:    g.acct.UserId,
			AccountId: g.acct.Id,
			RemoteId:  evt.Id,
			Timestamp: timestamp,
			Type:      typ,
			Resource:  repo,
			Origin:    from,
			Link:      link,
			Subject:   subject,
			Body:      title,
		}

	case "ForkEvent":
		typ := "fork"
		forkee := evt.Payload["forkee"].(map[string]interface{})
		forkName := forkee["full_name"].(string)
		owner := forkee["owner"].(map[string]interface{})
		from := owner["login"].(string)
		link := forkee["html_url"].(string)
		repo := evt.Repo.Name

		subject := fmt.Sprintf("Fork of %s created", repo)
		body := fmt.Sprintf("Repository %s was forked to %s.", repo, forkName)

		if !g.filter(typ, repo) {
			return
		}

		notf = &notification.Notification{
			UserId:    g.acct.UserId,
			AccountId: g.acct.Id,
			RemoteId:  evt.Id,
			Timestamp: timestamp,
			Type:      typ,
			Resource:  repo,
			Origin:    from,
			Link:      link,
			Subject:   subject,
			Body:      body,
		}

	case "CommitCommentEvent":
		typ := "commit_comment"
		comment := evt.Payload["comment"].(map[string]interface{})
		body := comment["body"].(string)
		user := comment["user"].(map[string]interface{})
		from := user["login"].(string)
		link := comment["html_url"].(string)
		repo := evt.Repo.Name

		if !g.filter(typ, repo) {
			return
		}

		subject := fmt.Sprintf("Commit comment in %s", repo)

		if len(body) > 140 {
			body = body[:140]
		}

		notf = &notification.Notification{
			UserId:    g.acct.UserId,
			AccountId: g.acct.Id,
			RemoteId:  evt.Id,
			Timestamp: timestamp,
			Type:      typ,
			Resource:  repo,
			Origin:    from,
			Link:      link,
			Subject:   subject,
			Body:      body,
		}

	case "PullRequestEvent", "PullRequestReviewCommentEvent":
		action := evt.Payload["action"].(string)
		pull := evt.Payload["pull_request"].(map[string]interface{})
		user := pull["user"].(map[string]interface{})
		from := user["login"].(string)
		title := pull["title"].(string)
		link := pull["html_url"].(string)
		repo := evt.Repo.Name

		if len(title) > 140 {
			title = title[:140]
		}

		var typ string
		var subject string

		if evt.Type == "PullRequestReviewCommentEvent" {
			if action != "created" {
				return
			}

			typ = "pull_comment"
			subject = fmt.Sprintf("New pull request comment in %s", repo)
		} else {
			switch action {
			case "assigned", "unassigned":
				if pull["assignee"].(string) != g.acct.Identity {
					return
				}
				typ = "pull_" + action
			case "opened":
				typ = "pull_opened"
			case "closed":
				typ = "pull_closed"
			case "reopened":
				typ = "pull_reopened"
			default:
				return
			}

			subject = fmt.Sprintf("Pull request %s in %s", action, repo)
		}

		if !g.filter(typ, repo) {
			return
		}

		notf = &notification.Notification{
			UserId:    g.acct.UserId,
			AccountId: g.acct.Id,
			RemoteId:  evt.Id,
			Timestamp: timestamp,
			Type:      typ,
			Resource:  repo,
			Origin:    from,
			Link:      link,
			Subject:   subject,
			Body:      title,
		}
	}

	return
}

func (g *gitHubBackend) sync() {
	lastNotf, err := notification.GetLastNotification(g.db,
		g.acct.UserId, g.acct.Id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("account.github: Failed to get last notification")
		return
	}
	g.lastNotf = lastNotf

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

	new := false

	for i, evt := range data {
		notf, done, e := g.parse(evt, i == 0)
		if e != nil {
			err = e
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("account.github: Failed to parse event")
		}

		if notf != nil {
			new = true

			err := notf.Initialize(g.db)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("account.github: Failed to parse event")
			}
		}

		if done {
			return
		}
	}

	if new {
		notification.PublishUpdate(g.acct.UserId)
	}
}

func (g *gitHubBackend) Run() {
	g.url = fmt.Sprintf("https://api.github.com/users/%s/received_events",
		g.acct.Identity)
	g.client = gitHubConf.NewClient(g.acct)
	g.interval = 0

	for {
		if g.interval > 0 {
			interval := time.Duration(utils.MinInt(g.interval, 180))
			time.Sleep(interval * time.Second)
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
