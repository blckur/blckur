package accounts

import (
	"encoding/base64"
	"fmt"
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/notification"
	"github.com/blckur/blckur/oauth"
	"github.com/blckur/blckur/settings"
	"github.com/dropbox/godropbox/errors"
	"labix.org/v2/mgo/bson"
	"strings"
	"time"
)

const (
	gmail = "gmail"
)

var (
	dateLayouts []string
	gmailConf   *oauth.Oauth2
)

func init() {
	dows := [...]string{"", "Mon, "}
	days := [...]string{"2", "02"}
	years := [...]string{"2006", "06"}
	seconds := [...]string{":05", ""}
	zones := [...]string{"-0700", "MST", "-0700 (MST)"}

	for _, dow := range dows {
		for _, day := range days {
			for _, year := range years {
				for _, second := range seconds {
					for _, zone := range zones {
						s := dow + day + " Jan " + year +
							" 15:04" + second + " " + zone
						dateLayouts = append(dateLayouts, s)
					}
				}
			}
		}
	}

	account.Register(gmail, "Gmail", Oauth2, GmailAuth{}, GmailClient{},
		[]*account.FilterType{
			&account.FilterType{
				Label: "All new messages",
				Type:  "all",
			},
			&account.FilterType{
				Label:     "Messages matching sender",
				Type:      "from",
				ValueType: "input",
				ValueLabel: "Enter complete or partial email address " +
					"of sender to match",
				ValueHolder: "Email address",
			},
			&account.FilterType{
				Label:       "Messages matching subject",
				Type:        "subject",
				ValueType:   "input",
				ValueLabel:  "Enter search term to match in email subject",
				ValueHolder: "Search term",
			},
			&account.FilterType{
				Label:       "Messages matching message body",
				Type:        "body",
				ValueType:   "input",
				ValueLabel:  "Enter search term to match in email body",
				ValueHolder: "Search term",
			},
			&account.FilterType{
				Label:     "Exclude messages matching sender",
				Type:      "not_from",
				ValueType: "input",
				ValueLabel: "Enter complete or partial email address " +
					"of sender to match",
				ValueHolder: "Email address",
			},
			&account.FilterType{
				Label:       "Exclude messages matching subject",
				Type:        "not_subject",
				ValueType:   "input",
				ValueLabel:  "Enter search term to match in email subject",
				ValueHolder: "Search term",
			},
			&account.FilterType{
				Label:       "Exclude messages matching message body",
				Type:        "not_body",
				ValueType:   "input",
				ValueLabel:  "Enter search term to match in email body",
				ValueHolder: "Search term",
			},
		}, func() {
			messenger.Register("settings", "google",
				func(_ *messenger.Message) {
					updateGmail()
				})
			updateGmail()
		})
}

func parseDate(date string) (parsed time.Time, err error) {
	for _, layout := range dateLayouts {
		parsed, err = time.Parse(layout, date)
		if err == nil {
			return
		}
	}

	err = &InvalidDate{
		errors.New("accounts.gmail: Failed to parse date"),
	}

	return
}

type gmailMessage struct {
	Id      string   `json:"id"`
	Labels  []string `json:"labelIds"`
	Snippet string   `json:"snippet"`
	Payload struct {
		Headers []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}
		Body struct {
			Size int    `json:"size"`
			Data string `json:"data"`
		} `json:"body"`
	} `json:"payload"`
}

type GmailClient struct {
	acct *account.Account
}

func (g *GmailClient) SetAccount(acct *account.Account) {
	g.acct = acct
}

func (g *GmailClient) Update(db *database.Database) (err error) {
	client := gmailConf.NewClient(g.acct)

	err = client.Refresh(db)
	if err != nil {
		return
	}

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
			var err error
			date, err = parseDate(header.Value)
			if err != nil {
				return
			}
		}
	}

	if msg.Id == lastNotf.RemoteId || date.Before(lastNotf.Timestamp) {
		done = true
		return
	}

	if force {
		notf = &notification.Notification{
			UserId:      g.acct.UserId,
			AccountId:   g.acct.Id,
			AccountType: gmail,
			RemoteId:    msg.Id,
			Timestamp:   date,
		}
	}

	if lastNotf == nil {
		done = true
		return
	}

	match := false
	body := ""

Loop:
	for _, filter := range g.acct.Filters {
		switch filter.Type {
		case "all":
			match = true
			break Loop
		case "from":
			match = strings.Contains(
				strings.ToLower(from),
				strings.ToLower(filter.Value))
			if match {
				break Loop
			}
		case "subject":
			match = strings.Contains(
				strings.ToLower(subject),
				strings.ToLower(filter.Value))
			if match {
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
			match = strings.Contains(body, strings.ToLower(filter.Value))
			if match {
				break Loop
			}
		}
	}

	if !match {
		return
	}

	limit := settings.Account.BodyLimit
	var bodySnippet string
	if len(msg.Snippet) > limit {
		bodySnippet = msg.Snippet[:limit] + "..."
	} else {
		bodySnippet = msg.Snippet
	}

	notf = &notification.Notification{
		UserId:      g.acct.UserId,
		AccountId:   g.acct.Id,
		AccountType: gmail,
		RemoteId:    msg.Id,
		Timestamp:   date,
		Type:        "email",
		Origin:      from,
		Link:        "https://mail.google.com/mail/u/0/#inbox/" + msg.Id,
		Subject:     subject,
		Body:        bodySnippet,
	}

	return
}

func (g *GmailClient) Sync(db *database.Database) (err error) {
	client := gmailConf.NewClient(g.acct)

	err = client.Refresh(db)
	if err != nil {
		return
	}

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
		messages := struct {
			Messages []struct {
				Id string `json:"id"`
			} `json:"messages"`
			NextPageToken string `json:"nextPageToken"`
		}{}

		url := fmt.Sprintf("https://www.googleapis.com/gmail/v1"+
			"/users/me/messages?labelIds=INBOX&maxResults=%d"+
			"&includeSpamTrash=false", msgCount)

		if pageToken != "" {
			url += "&pageToken=" + pageToken
		}

		err = client.GetJson(url, &messages)
		if err != nil {
			return
		}

		pageToken = messages.NextPageToken

		for j, msg := range messages.Messages {
			data := &gmailMessage{}

			url := fmt.Sprintf("https://www.googleapis.com/gmail/v1"+
				"/users/me/messages/%s?format=full", msg.Id)

			err = client.GetJson(url, data)
			if err != nil {
				return
			}

			notf, done := g.parseMessage(data, lastNotf, i == 0 && j == 0)
			if notf != nil {
				notfs = append(notfs, notf)
			}

			if done {
				break Loop
			}
		}
	}

	pub := notification.NewPublisher(g.acct.UserId.Hex())
	defer pub.Close()

	for i := len(notfs) - 1; i >= 0; i-- {
		notf := notfs[i]

		err = notf.Initialize(db)
		if err != nil {
			return
		}

		if notf.Type != "" {
			err = pub.New(notf)
			if err != nil {
				return
			}
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
	code string) (acct *account.Account, err error) {

	coll := db.Accounts()

	auth, err := gmailConf.Authorize(db, state, code)
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

func updateGmail() {
	gmailConf = &oauth.Oauth2{
		Type:         gmail,
		ClientId:     settings.Google.ClientId,
		ClientSecret: settings.Google.ClientSecret,
		CallbackUrl:  settings.System.Domain + "/callback/gmail",
		AuthUrl:      "https://accounts.google.com/o/oauth2/auth",
		TokenUrl:     "https://www.googleapis.com/oauth2/v3/token",
		Scopes: []string{
			"https://www.googleapis.com/auth/gmail.readonly",
		},
	}
	gmailConf.Config()
}
