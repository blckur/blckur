package account

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/notification"
	"github.com/blckur/blckur/oauth"
	"github.com/dropbox/godropbox/container/set"
	"labix.org/v2/mgo/bson"
	"encoding/base64"
	"fmt"
	"time"
	"strings"
)

var (
	gmailConf *oauth.Oauth2
)

func init() {
	register("gmail", "Gmail", OAUTH2, GmailAuth{}, GmailClient{}, func() {
		messenger.Register("settings", "google", func(_ *messenger.Message) {
			updateGmail()
		})
		updateGmail()
	})
}

type gmailMessage struct {
	Id string `json:"id"`
	Labels []string `json:"labelIds"`
	Snippet string `json:"snippet"`
	Payload struct {
		Headers []struct {
			Name string `json:"name"`
			Value string `json:"value"`
		}
		Body struct {
			Size int `json:"size"`
			Data string `json:"data"`
		} `json:"body"`
	} `json:"payload"`
}

type GmailClient struct {
	acct *Account
}

func (g *GmailClient) setAccount(acct *Account) {
	g.acct = acct
}

func (g *GmailClient) newClient() (client *oauth.Oauth2Client) {
	client = gmailConf.NewClient(g.acct.UserId, g.acct.Oauth2AccTokn,
		g.acct.Oauth2RefTokn, g.acct.Oauth2Exp)
	return
}

func (g *GmailClient) refresh(db *database.Database,
		client *oauth.Oauth2Client) (err error) {
	refreshed, err := client.Check()
	if err != nil {
		return
	}

	if refreshed {
		coll := db.Accounts()

		g.acct.Oauth2AccTokn = client.AccessToken
		g.acct.Oauth2RefTokn = client.RefreshToken
		g.acct.Oauth2Exp = client.Expiry

		fields := set.NewSet("oauth2_acc_tokn", "oauth2_ref_tokn",
			"oauth2_exp")

		err = coll.CommitFields(g.acct.Id, g, fields)
		if err != nil {
			err = database.ParseError(err)
			return
		}
	}

	return
}

func (g *GmailClient) Update(db *database.Database) (err error) {
	client := g.newClient()
	g.refresh(db, client)

	data := struct {
		EmailAddress string `json:"emailAddress"`
	}{}

	err = client.GetJson(
		"https://www.googleapis.com/gmail/v1/users/me/profile", &data)
	if err != nil {
		return
	}

	g.acct.Identity = data.EmailAddress

	return
}

func (g *GmailClient) parseMessage(msg *gmailMessage,
		lastNotf *notification.Notification, force bool) (
		notf *notification.Notification, done bool) {
	done = false
	from := ""
	subject := ""
	var date time.Time

	for _, header := range msg.Payload.Headers {
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
		return
	}

	if force || lastNotf == nil {
		notf = &notification.Notification{
			UserId: g.acct.UserId,
			AccountId: g.acct.Id,
			RemoteId: msg.Id,
			Timestamp: date,
		}
	}

	if lastNotf == nil {
		done = true
		return
	}

	if date.Before(lastNotf.Timestamp) || msg.Id == lastNotf.RemoteId {
		done = true
		return
	}

	match := false
	matchType := ""
	body := ""

	Loop:
	for _, alrt := range g.acct.Alerts {
		switch (alrt.Type) {
		case "all":
			match = true
			matchType = "all"
			break Loop
		case "from":
			match = strings.Contains(
				strings.ToLower(from),
				strings.ToLower(alrt.Value))
			if match {
				matchType = "from"
				break Loop
			}
		case "subject":
			match = strings.Contains(
				strings.ToLower(subject),
				strings.ToLower(alrt.Value))
			if match {
				matchType = "subject"
				break Loop
			}
		case "body":
			if body == "" {
				bodyByt, err := base64.URLEncoding.DecodeString(
					msg.Payload.Body.Data)
				if err != nil {
					body = "-"
				} else {
					body = strings.ToLower(string(bodyByt))
				}
			}
			match = strings.Contains(body, strings.ToLower(alrt.Value))
			if match {
				matchType = "body"
				break Loop
			}
		}
	}

	if !match {
		return
	}

	var bodySnippet string
	if len(msg.Snippet) > 140 {
		bodySnippet = msg.Snippet[:140]
	} else {
		bodySnippet = msg.Snippet
	}

	notf = &notification.Notification{
		UserId: g.acct.UserId,
		AccountId: g.acct.Id,
		RemoteId: msg.Id,
		Timestamp: date,
		Type: matchType,
		Origin: from,
		Subject: subject,
		Body: bodySnippet,
	}

	return
}

func (g *GmailClient) Sync(db *database.Database) (err error) {
	client := g.newClient()
	g.refresh(db, client)

	lastNotf, err := notification.GetLastNotification(db,
		g.acct.UserId, g.acct.Id)
	if err != nil {
		return
	}

	var msgCount int
	if lastNotf == nil {
		msgCount = 3
	} else {
		msgCount = 10
	}

	pageToken := ""

	notfs := []*notification.Notification{}

	n := settings.Google.MaxMsg / 10
	Loop:
	for i := 0; i < n; i++ {
		messages := struct{
			Messages []struct{
				Id string `json:"id"`
			} `json:"messages"`
			NextPageToken string `json:"nextPageToken"`
		}{}

		url := fmt.Sprintf("https://www.googleapis.com/gmail/v1" +
			"/users/me/messages?labelIds=INBOX&maxResults=%d" +
			"&includeSpamTrash=false", msgCount)

		if pageToken != "" {
			url += "&pageToken=" + pageToken
		}

		err = client.GetJson(url, &messages)
		if err != nil {
			return
		}

		pageToken = messages.NextPageToken

		for i, msg := range messages.Messages {
			data := &gmailMessage{}

			url := fmt.Sprintf("https://www.googleapis.com/gmail/v1" +
				"/users/me/messages/%s?format=full", msg.Id)

			err = client.GetJson(url, data)
			if err != nil {
				return
			}

			notf, done := g.parseMessage(data, lastNotf, i == 0)
			if notf != nil {
				notfs = append(notfs, notf)
			}

			if done {
				break Loop
			}
		}
	}

	for i := len(notfs) - 1; i >= 0; i-- {
		notf := notfs[i]

		err = notf.Initialize(db)
		if err != nil {
			return
		}
	}

	return
}

type GmailAuth struct {
}

func (g *GmailAuth) Request(db *database.Database, userId bson.ObjectId) (
		url string, err error) {
	url, err = gmailConf.Request(db, userId)
	if err != nil {
		return
	}

	return
}

func (g *GmailAuth) Authorize(db *database.Database, state string,
		code string) (acct *Account, err error) {
	coll := db.Accounts()

	auth, err := gmailConf.Authorize(db, state, code)
	if err != nil {
		return
	}

	acct = &Account{
		UserId: auth.UserId,
		Type: "gmail",
		Oauth2AccTokn: auth.AccessToken,
		Oauth2RefTokn: auth.RefreshToken,
		Oauth2Exp: auth.Expiry,
		coll: coll,
	}

	client, err := acct.GetClient()
	if err != nil {
		return
	}

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
