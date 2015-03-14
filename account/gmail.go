package account

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/notification"
	"github.com/blckur/blckur/oauth"
	"github.com/dropbox/godropbox/container/set"
	"labix.org/v2/mgo/bson"
	"fmt"
	"time"
)

var (
	gmailConf *oauth.Oauth2
)

type Gmail struct {
	Account `bson:",inline"`
}

func (g *Gmail) NewClient() (client *oauth.Oauth2Client) {
	client = gmailConf.NewClient(g.UserId, g.Oauth2AccTokn,
		g.Oauth2RefTokn, g.Oauth2Exp)
	return
}

func (g *Gmail) Refresh(db *database.Database, client *oauth.Oauth2Client) (
		err error) {
	refreshed, err := client.Check()
	if err != nil {
		return
	}

	if refreshed {
		coll := db.Accounts()

		g.Oauth2AccTokn = client.AccessToken
		g.Oauth2RefTokn = client.RefreshToken
		g.Oauth2Exp = client.Expiry

		fields := set.NewSet("oauth2_acc_tokn", "oauth2_ref_tokn",
			"oauth2_exp")

		err = coll.CommitFields(g.Id, g, fields)
		if err != nil {
			err = database.ParseError(err)
			return
		}
	}

	return
}


func (g *Gmail) Update(db *database.Database) (err error) {
	client := g.NewClient()
	g.Refresh(db, client)

	data := struct {
		EmailAddress string `json:"emailAddress"`
	}{}

	err = client.GetJson(
		"https://www.googleapis.com/gmail/v1/users/me/profile", &data)
	if err != nil {
		return
	}

	g.Identity = data.EmailAddress

	return
}

func (g *Gmail) Sync(db *database.Database) (err error) {
	client := g.NewClient()
	g.Refresh(db, client)

	lastNotf, err := notification.GetLastNotification(db, g.UserId, g.Type)
	if err != nil {
		return
	}

	var msgCount int
	if lastNotf == nil {
		msgCount = 1
	} else {
		msgCount = 10
	}

	messages := struct{
		Messages []struct{
			Id string `json:"id"`
		} `json:"messages"`
	}{}

	url := fmt.Sprintf("https://www.googleapis.com/gmail/v1" +
		"/users/me/messages?labelIds=INBOX&fields=messages(id)" +
		"&maxResults=%d&includeSpamTrash=false", msgCount)

	err = client.GetJson(url, &messages)
	if err != nil {
		return
	}

	for _, msg := range messages.Messages {
		data := struct{
			Id string `json:"id"`
			Labels []string `json:"labelIds"`
			Snippet string `json:"snippet"`
			Payload struct{
				Headers []struct{
					Name string `json:"name"`
					Value string `json:"value"`
				}
			}
			Body struct{
				Size int `json:"size"`
				Data string `json:"data"`
			}
		}{}

		url := fmt.Sprintf("https://www.googleapis.com/gmail/v1" +
			"/users/me/messages/%s?format=full", msg.Id)

		err = client.GetJson(url, &data)
		if err != nil {
			return
		}

		from := ""
		subject := ""
		var date time.Time

		for _, header := range data.Payload.Headers {
			if header.Name == "From" {
				from = header.Value
			} else if header.Name == "Subject" {
				subject = header.Value
			} else if header.Name == "Date" {
				date, _ = time.Parse("Mon, 02 Jan 2006 15:04:05 -0700",
					header.Value)

			}
		}

		if date.IsZero() {
			continue
		}

		var notf *notification.Notification
		if lastNotf == nil {
			notf = &notification.Notification{
				UserId: g.UserId,
				RemoteId: data.Id,
				Timestamp: date,
				AccountType: g.Type,
			}
			break
		} else {
			if date.Before(lastNotf.Timestamp) ||
					data.Id == lastNotf.RemoteId {
				break
			}

			notf = &notification.Notification{
				UserId: g.UserId,
				RemoteId: data.Id,
				Timestamp: date,
				AccountType: g.Type,
				Type: "all",
				Origin: from,
				Subject: subject,
				Body: data.Snippet,
			}
		}

		err = notf.Initialize(db)
		if err != nil {
			return
		}
	}

	return
}

func ReqGmail(db *database.Database, userId bson.ObjectId) (
		url string, err error) {
	url, err = gmailConf.Request(db, userId)
	if err != nil {
		return
	}

	return
}

func AuthGmail(db *database.Database, state string, code string) (
		acct *Account, err error) {
	coll := db.Accounts()

	client, err := gmailConf.Authorize(db, state, code)
	if err != nil {
		return
	}

	data := struct {
		EmailAddress string `bson:"emailAddress"`
	}{}

	err = client.GetJson(
		"https://www.googleapis.com/gmail/v1/users/me/profile", &data)
	if err != nil {
		return
	}

	acct = &Account{
		UserId: client.UserId,
		Type: "gmail",
		Identity: data.EmailAddress,
		Oauth2AccTokn: client.AccessToken,
		Oauth2RefTokn: client.RefreshToken,
		Oauth2Exp: client.Expiry,
		coll: coll,
	}

	err = coll.Insert(acct)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func updateGmail() {
	gmailConf = &oauth.Oauth2{
		Type: "gmail",
		ClientId: settings.Google.ClientId,
		ClientSecret: settings.Google.ClientSecret,
		CallbackUrl: settings.System.Domain + "/callback/gmail",
		AuthUrl: "https://accounts.google.com/o/oauth2/auth",
		TokenUrl: "https://www.googleapis.com/oauth2/v3/token",
		Scopes: []string{
			"https://www.googleapis.com/auth/gmail.readonly",
		},
	}
	gmailConf.Config()
}

func InitGmail() {
	requires.After("settings")
	requires.Before("messenger")

	messenger.Register("settings", "google", func(_ *messenger.Message) {
		updateGmail()
	})
	updateGmail()

	requires.Register("account")
}
