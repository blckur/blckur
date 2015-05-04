package accounts

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/cache"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/notification"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/stream"
	"github.com/dropbox/godropbox/errors"
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

const (
	hackerNews = "hackernews"
)

func init() {
	account.Register(hackerNews, "Hacker News", Oauth2,
		HackerNewsAuth{}, HackerNewsClient{},
		[]*account.FilterType{
			&account.FilterType{
				Label: "All new stories",
				Type:  "all",
			},
			&account.FilterType{
				Label:       "Story matching title",
				Type:        "story_title",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial text to match",
				ValueHolder: "Title text",
			},
			&account.FilterType{
				Label:       "Story matching url",
				Type:        "story_url",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial url to match",
				ValueHolder: "Url to match",
			},
			&account.FilterType{
				Label:       "Story matching text",
				Type:        "story_text",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial text to match",
				ValueHolder: "Story text",
			},
		}, nil)
}

type HackerNewsClient struct {
	acct *account.Account
}

func (c *HackerNewsClient) SetAccount(acct *account.Account) {
	c.acct = acct
}

func (c *HackerNewsClient) Update(_ *database.Database) (err error) {
	return
}

func (c *HackerNewsClient) Sync(db *database.Database) (err error) {
	backend := &hackerNewsBackend{
		db:   db,
		acct: c.acct,
	}
	stream := stream.NewStream(db, c.acct.Id, c.acct.UserId, backend)

	err = stream.Start()
	if err != nil {
		return
	}

	return
}

type hackerNewsBackend struct {
	db   *database.Database
	acct *account.Account
	stop bool
}

type hackerNewsStory struct {
	Id     int       `json:"id"`
	Title  string    `json:"title"`
	By     string    `json:"by"`
	Time   time.Time `json:"time"`
	Text   string    `json:"text"`
	Url    string    `json:"url"`
	Domain string    `json:"domain"`
}

func (b *hackerNewsBackend) parse(story *hackerNewsStory) (err error) {
	match := false

	for _, filter := range b.acct.Filters {
		switch filter.Type {
		case "all":
			match = true
		}
	}

	if !match {
		return
	}

	var body string
	if story.Domain == "" {
		body = fmt.Sprintf("New story by %s", story.By)
	} else {
		body = fmt.Sprintf("New story on %s by %s", story.Domain, story.By)
	}

	notf := &notification.Notification{
		UserId:      b.acct.UserId,
		AccountId:   b.acct.Id,
		AccountType: hackerNews,
		RemoteId:    strconv.Itoa(story.Id),
		Timestamp:   story.Time,
		Type:        "story",
		Origin:      story.By,
		Link:        story.Url,
		Subject:     story.Title,
		Body:        body,
	}

	newNotf, err := notf.Initialize(b.db)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("account.hackernews: Failed to parse event")
	} else if newNotf {
		pub := notification.NewPublisher(b.acct.UserId.Hex())
		defer pub.Close()

		err = pub.New(notf)
		if err != nil {
			return
		}
	}

	return
}

func (b *hackerNewsBackend) backlog() (err error) {
	conn := cache.Get()

	storyIdsStr, err := conn.GetString("new_stories")
	if err != nil {
		return
	}

	if storyIdsStr == "" {
		return
	}

	storyIds := []int{}

	err = json.Unmarshal([]byte(storyIdsStr), &storyIds)
	if err != nil {
		err = &ParseError{
			errors.Wrap(err, "Hacker News backlog parse error"),
		}
		return
	}

	for _, storyId := range storyIds {
		storyStr, e := conn.GetString(fmt.Sprintf("story-%d", storyId))
		if e != nil {
			err = e
			return
		}

		story := &hackerNewsStory{}

		e = json.Unmarshal([]byte(storyStr), story)
		if e != nil {
			err = e
			return
		}

		err = b.parse(story)
		if err != nil {
			return
		}
	}

	return
}

func (b *hackerNewsBackend) Run() {
	lst := cache.Subscribe("hacker_news")
	defer lst.Close()

	sub := lst.Listen()

	go func() {
		err := b.backlog()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("accounts.hackernews: Unknown backlog error")
		}
	}()

	for msg := range sub {
		if msg == nil {
			return
		}

		go func() {
			story := &hackerNewsStory{}

			err := json.Unmarshal([]byte(msg.Data.(string)), story)
			if err != nil {
				err = &ParseError{
					errors.Wrap(err, "Hacker News sub error"),
				}
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("accounts.hackernews: Unknown error")
			}

			err = b.parse(story)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("accounts.hackernews: Unknown parse error")
			}
		}()
	}
}

func (b *hackerNewsBackend) Stop() {
	b.stop = true
}

type HackerNewsAuth struct{}

func (a *HackerNewsAuth) Request(db *database.Database,
	userId bson.ObjectId) (url string, err error) {

	coll := db.Accounts()

	acct := &account.Account{
		UserId: userId,
		Type:   hackerNews,
		New:    true,
	}

	err = coll.Insert(acct)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	url = settings.System.AppHome

	return
}

func (a *HackerNewsAuth) Authorize(_ *database.Database, _ string,
	_ string) (acct *account.Account, err error) {

	return
}
