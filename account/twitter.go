package account

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/oauth"
	"github.com/blckur/blckur/notification"
	"github.com/blckur/blckur/stream"
	"github.com/ChimeraCoder/anaconda"
	"github.com/Sirupsen/logrus"
	"labix.org/v2/mgo/bson"
	"crypto/md5"
	"net/url"
	"time"
	"strconv"
	"encoding/hex"
)

var (
	twitterConf *oauth.Oauth1
)

type TwitterClient struct {
	acct *Account
}

func init() {
	register("twitter", OAUTH1, TwitterClient{}, TwitterAuth{}, func() {
		messenger.Register("settings", "twitter", func(_ *messenger.Message) {
			updateTwitter()
		})
		updateTwitter()
	})
}

func (t *TwitterClient) setAccount(acct *Account) {
	t.acct = acct
}

func (t *TwitterClient) newClient() (client *oauth.Oauth1Client) {
	client = twitterConf.NewClient(t.acct.UserId, t.acct.OauthTokn,
		t.acct.OauthSec)
	return
}

func (t *TwitterClient) Update(db *database.Database) (err error) {
	client := t.newClient()

	data := struct {
		IdStr string `json:"id_str"`
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
		db: db,
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
	db *database.Database
	stream *anaconda.Stream
	acct *Account
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
		_ = evt
		return
	} else if evt, ok := evtInf.(anaconda.EventTweet); ok {
		if evt.Target.IdStr != b.acct.IdentityId ||
				evt.Source.IdStr == b.acct.IdentityId {
			return
		}

		var subject string
		switch (evt.Event.Event) {
		case "favorite":
			subject = "Tweet favorited by "
		case "unfavorite":
			subject = "Tweet unfavorited by "
		default:
			return
		}

		origin := "@" + evt.Source.ScreenName
		timestamp = evt.CreatedAt
		subject += origin

		notf = &notification.Notification{
			UserId: b.acct.UserId,
			AccountId: b.acct.Id,
			Type: evt.Event.Event,
			Resource: evt.TargetObject.IdStr,
			Origin: origin,
			Subject: subject,
			Body: evt.TargetObject.Text,
		}
	} else if evt, ok := evtInf.(anaconda.Event); ok {
		if evt.Target.IdStr != b.acct.IdentityId || evt.Event != "follow" {
			return
		}

		origin := "@" + evt.Source.ScreenName
		timestamp = evt.CreatedAt
		notf = &notification.Notification{
			UserId: b.acct.UserId,
			AccountId: b.acct.Id,
			Type: evt.Event,
			Resource: evt.Target.IdStr,
			Origin: origin,
			Subject: "New follower " + origin,
		}
	} else if evt, ok := evtInf.(anaconda.StatusDeletionNotice); ok {
		coll := b.db.Notifications()

		err = coll.Remove(bson.M{
			"account_id": b.acct.Id,
			"resource": evt.IdStr,
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

	return
}

func hashEvent(id string, timestamp time.Time) (hashStr string) {
	hash := md5.Sum([]byte(id + strconv.FormatInt(timestamp.Unix(), 10)))
	hashStr = hex.EncodeToString(hash[:])
	return
}

type TwitterAuth struct {
}

func (t *TwitterAuth) Request(db *database.Database, userId bson.ObjectId) (
		url string, err error) {
	url, err = twitterConf.Request(db, userId)
	if err != nil {
		return
	}

	return
}

func (t *TwitterAuth) Authorize(db *database.Database, token string,
		code string) (acct *Account, err error) {
	coll := db.Accounts()

	auth, err := twitterConf.Authorize(db, token, code)
	if err != nil {
		return
	}

	acct = &Account{
		UserId: auth.UserId,
		Type: "twitter",
		OauthTokn: auth.Token,
		OauthSec: auth.Secret,
		coll: coll,
	}

	client := acct.GetClient()
	err = client.Update(db)
	if err != nil {
		return
	}

	err = coll.Insert(acct)
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
		Type: "twitter",
		ConsumerKey: settings.Twitter.ConsumerKey,
		ConsumerSecret: settings.Twitter.ConsumerSecret,
		ReqTokenUrl: "https://api.twitter.com/oauth/request_token",
		AuthTokenUrl: "https://api.twitter.com/oauth/authorize",
		AccsTokenUrl: "https://api.twitter.com/oauth/access_token",
		CallbackUrl: settings.System.Domain + "/callback/twitter",
	}
	conf.Config()

	twitterConf = conf
}
