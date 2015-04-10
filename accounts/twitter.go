package accounts

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/Sirupsen/logrus"
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/notification"
	"github.com/blckur/blckur/oauth"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/stream"
	"labix.org/v2/mgo/bson"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	TWITTER = "twitter"
)

var (
	twitterConf *oauth.Oauth1
)

type TwitterClient struct {
	acct *account.Account
}

func init() {
	account.Register(TWITTER, "Twitter", OAUTH1,
		TwitterAuth{}, TwitterClient{},
		[]*account.FilterType{
			&account.FilterType{
				Label: "All new events",
				Type:  "all",
			},
			&account.FilterType{
				Label:       "Events matching user",
				Type:        "from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial twitter handle",
				ValueHolder: "Twitter handle",
			},
			&account.FilterType{
				Label: "All new followers",
				Type:  "follow_all",
			},
			&account.FilterType{
				Label:       "Followers matching user",
				Type:        "follow_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial twitter handle",
				ValueHolder: "Twitter handle",
			},
			&account.FilterType{
				Label: "All new favorited tweets",
				Type:  "favorite_all",
			},
			&account.FilterType{
				Label:       "Favorited by matching user",
				Type:        "favorite_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial twitter handle",
				ValueHolder: "Twitter handle",
			},
			&account.FilterType{
				Label: "All new unfavorited tweets",
				Type:  "unfavorite_all",
			},
			&account.FilterType{
				Label:       "Unfavorited by matching user",
				Type:        "unfavorite_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial twitter handle",
				ValueHolder: "Twitter handle",
			},
			&account.FilterType{
				Label: "All new retweets",
				Type:  "retweet_all",
			},
			&account.FilterType{
				Label:       "Retweets by matching user",
				Type:        "retweet_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial twitter handle",
				ValueHolder: "Twitter handle",
			},
			&account.FilterType{
				Label: "All new tweet replies",
				Type:  "reply_all",
			},
			&account.FilterType{
				Label:       "Tweet replies by matching user",
				Type:        "reply_from",
				ValueType:   "input",
				ValueLabel:  "Enter complete or partial twitter handle",
				ValueHolder: "Twitter handle",
			},
		}, func() {
			messenger.Register("settings", TWITTER,
				func(_ *messenger.Message) {
					updateTwitter()
				})
			updateTwitter()
		})
}

func (t *TwitterClient) SetAccount(acct *account.Account) {
	t.acct = acct
}

func (t *TwitterClient) Update(db *database.Database) (err error) {
	client := twitterConf.NewClient(t.acct)

	data := struct {
		IdStr      string `json:"id_str"`
		ScreenName string `json:"screen_name"`
	}{}

	err = client.GetJson(
		"https://api.twitter.com/1.1/account/verify_credentials.json",
		nil, &data)
	if err != nil {
		return
	}

	t.acct.Identity = "@" + data.ScreenName
	t.acct.IdentityId = data.IdStr

	return
}

func (t *TwitterClient) Sync(db *database.Database) (err error) {
	backend := &twitterBackend{
		db:   db,
		acct: t.acct,
	}
	stream := stream.NewStream(db, t.acct.Id, backend)

	err = stream.Start()
	if err != nil {
		return
	}

	return
}

type twitterBackend struct {
	db     *database.Database
	stream *anaconda.Stream
	acct   *account.Account
}

func (b *twitterBackend) filter(typ string, origin string) bool {
	for _, filter := range b.acct.Filters {
		switch filter.Type {
		case "all", typ + "_all":
			return true
		case typ + "_from":
			if strings.Contains(origin, filter.Value) {
				return true
			}
		}
	}

	return false
}

func (b *twitterBackend) newClient() (client *anaconda.TwitterApi) {
	client = anaconda.NewTwitterApi(b.acct.OauthTokn, b.acct.OauthSec)
	return
}

func (b *twitterBackend) Run() {
	client := b.newClient()

	streamVals := url.Values{}
	streamVals.Add("with", "user")
	streamVals.Add("replies", "all")

	s := client.UserStream(streamVals)
	b.stream = &s

	for obj := range b.stream.C {
		_, err := b.handle(obj)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("account.twitter: Stream handle error")
		}
	}
}

func (b *twitterBackend) Stop() {
	b.stream.Interrupt()
	b.stream.End()
}

func (b *twitterBackend) handle(evtInf interface{}) (
	notf *notification.Notification, err error) {
	var timestamp string

	if evt, ok := evtInf.(anaconda.Tweet); ok {
		var evtType string
		var subject string

		if evt.RetweetedStatus != nil {
			evtType = "retweet"
			subject = "Tweet retweeted by "
		} else if evt.InReplyToUserIdStr == b.acct.IdentityId {
			if evt.InReplyToStatusIdStr == "" {
				evtType = "mention"
				subject = "Mentioned by "
			} else {
				evtType = "reply"
				subject = "Tweet reply by "
			}
		} else {
			return
		}

		origin := "@" + evt.User.ScreenName
		timestamp = evt.CreatedAt
		subject += fmt.Sprintf("%s (%s)", evt.User.Name, origin)

		if !b.filter(evtType, origin) {
			return
		}

		notf = &notification.Notification{
			UserId:      b.acct.UserId,
			AccountId:   b.acct.Id,
			AccountType: TWITTER,
			Type:        evtType,
			Resource:    evt.IdStr,
			Origin:      origin,
			Subject:     subject,
			Body:        evt.Text,
		}
	} else if evt, ok := evtInf.(anaconda.EventTweet); ok {
		if evt.Target.IdStr != b.acct.IdentityId ||
			evt.Source.IdStr == b.acct.IdentityId {
			return
		}

		var subject string
		switch evt.Event.Event {
		case "favorite":
			subject = "Tweet favorited by "
		case "unfavorite":
			subject = "Tweet unfavorited by "
		default:
			return
		}

		origin := "@" + evt.Source.ScreenName
		timestamp = evt.CreatedAt
		subject += fmt.Sprintf("%s (%s)", evt.Source.Name, origin)

		if !b.filter(evt.Event.Event, origin) {
			return
		}

		notf = &notification.Notification{
			UserId:      b.acct.UserId,
			AccountId:   b.acct.Id,
			AccountType: TWITTER,
			Type:        evt.Event.Event,
			Resource:    evt.TargetObject.IdStr,
			Origin:      origin,
			Subject:     subject,
			Body:        evt.TargetObject.Text,
		}
	} else if evt, ok := evtInf.(anaconda.Event); ok {
		if evt.Target.IdStr != b.acct.IdentityId || evt.Event != "follow" {
			return
		}

		origin := "@" + evt.Source.ScreenName
		timestamp = evt.CreatedAt

		if !b.filter(evt.Event, origin) {
			return
		}

		notf = &notification.Notification{
			UserId:      b.acct.UserId,
			AccountId:   b.acct.Id,
			AccountType: TWITTER,
			Type:        evt.Event,
			Resource:    evt.Target.IdStr,
			Origin:      origin,
			Subject:     "New follower " + origin,
		}
	} else if evt, ok := evtInf.(anaconda.StatusDeletionNotice); ok {
		coll := b.db.Notifications()

		err = coll.Remove(bson.M{
			"account_id": b.acct.Id,
			"resource":   evt.IdStr,
		})
		if err != nil {
			err = database.ParseError(err)
			return
		}

		return
	} else {
		return
	}

	notf.Timestamp, err = time.Parse(
		"Mon Jan 02 15:04:05 -0700 2006", timestamp)
	if err != nil {
		notf = nil
		err = nil
	}

	notf.RemoteId = hashEvent(b.acct.Id.Hex(), notf.Timestamp)

	err = notf.Initialize(b.db)
	if err != nil {
		return
	}

	notification.PublishUpdate(b.acct.UserId)

	return
}

func hashEvent(id string, timestamp time.Time) (hashStr string) {
	hash := md5.Sum([]byte(id + strconv.FormatInt(timestamp.Unix(), 10)))
	hashStr = hex.EncodeToString(hash[:])
	return
}

type TwitterAuth struct{}

func (t *TwitterAuth) Request(db *database.Database, userId bson.ObjectId) (
	url string, err error) {

	url, err = twitterConf.Request(db, userId)
	if err != nil {
		return
	}

	return
}

func (t *TwitterAuth) Authorize(db *database.Database, token string,
	code string) (acct *account.Account, err error) {

	coll := db.Accounts()

	auth, err := twitterConf.Authorize(db, token, code)
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

func updateTwitter() {
	anaconda.SetConsumerKey(settings.Twitter.ConsumerKey)
	anaconda.SetConsumerSecret(settings.Twitter.ConsumerSecret)

	conf := &oauth.Oauth1{
		Type:           TWITTER,
		ConsumerKey:    settings.Twitter.ConsumerKey,
		ConsumerSecret: settings.Twitter.ConsumerSecret,
		CallbackUrl:    settings.System.Domain + "/callback/twitter",
		ReqTokenUrl:    "https://api.twitter.com/oauth/request_token",
		AuthTokenUrl:   "https://api.twitter.com/oauth/authorize",
		AccsTokenUrl:   "https://api.twitter.com/oauth/access_token",
	}
	conf.Config()

	twitterConf = conf
}
