package accounts

import (
	"fmt"
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/notification"
	"github.com/blckur/blckur/oauth"
	"github.com/blckur/blckur/settings"
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

const (
	digitalOcean = "digitalocean"
)

var (
	digitalOceanConf *oauth.Oauth2
)

func init() {
	account.Register(digitalOcean, "DigitalOcean", Oauth2,
		DigitalOceanAuth{}, DigitalOceanClient{},
		[]*account.FilterType{
			&account.FilterType{
				Label: "All new events",
				Type:  "all",
			},
			&account.FilterType{
				Label: "Droplet created",
				Type:  "create",
			},
			&account.FilterType{
				Label: "Droplet destroyed",
				Type:  "destroy",
			},
			&account.FilterType{
				Label: "Droplet password reset",
				Type:  "password_reset",
			},
		}, func() {
			messenger.Register("settings", digitalOcean,
				func(_ *messenger.Message) {
					updateDigitalOcean()
				})
			updateDigitalOcean()
		})
}

type DigitalOceanClient struct {
	acct *account.Account
}

func (d *DigitalOceanClient) SetAccount(acct *account.Account) {
	d.acct = acct
}

func (d *DigitalOceanClient) Update(db *database.Database) (err error) {
	client := digitalOceanConf.NewClient(d.acct)

	err = client.Refresh(db)
	if err != nil {
		return
	}

	data := struct {
		Account struct {
			Email string `json:"email"`
		} `json:"account"`
	}{}

	err = client.GetJson(
		"https://api.digitalocean.com/v2/account", &data)
	if err != nil {
		return
	}

	d.acct.Identity = data.Account.Email

	return
}

func (d *DigitalOceanClient) filter(typ string) bool {
	for _, filter := range d.acct.Filters {
		if filter.Type == "all" || filter.Type == typ {
			return true
		}
	}

	return false
}

func (d *DigitalOceanClient) Sync(db *database.Database) (err error) {
	client := digitalOceanConf.NewClient(d.acct)

	err = client.Refresh(db)
	if err != nil {
		return
	}

	lastNotf, err := notification.GetLastNotification(db,
		d.acct.UserId, d.acct.Id)
	if err != nil {
		return
	}

	var msgCount int
	if lastNotf == nil {
		msgCount = 3
	} else {
		msgCount = 10
	}

	url := fmt.Sprintf("https://api.digitalocean.com/v2"+
		"/actions?page=1&per_page=%d", msgCount)

	notfs := []*notification.Notification{}

	n := settings.Google.MaxMsg / 10
Loop:
	for i := 0; i < n; i++ {
		actions := struct {
			Actions []struct {
				Id        int    `json:"id"`
				StartedAt string `json:"started_at"`
				Type      string `json:"type"`
			} `json:"actions"`
			Links struct {
				Pages struct {
					Next string `json:"next"`
				} `json:"pages"`
			} `json:"links"`
		}{}

		err = client.GetJson(url, &actions)
		if err != nil {
			return
		}

		url = actions.Links.Pages.Next

		for _, action := range actions.Actions {
			timestamp, _ := time.Parse("2006-01-02T15:04:05Z",
				action.StartedAt)
			if timestamp.IsZero() {
				continue
			}

			var subject string
			switch action.Type {
			case "create":
				subject = "Droplet created"
			case "destroy":
				subject = "Droplet destroyed"
			case "password_reset":
				subject = "Droplet password reset"
			default:
				continue
			}

			if !d.filter(action.Type) {
				continue
			}

			if lastNotf == nil {
				notf := &notification.Notification{
					UserId:      d.acct.UserId,
					AccountId:   d.acct.Id,
					AccountType: digitalOcean,
					RemoteId:    strconv.Itoa(action.Id),
					Timestamp:   timestamp,
				}
				notfs = append(notfs, notf)
				break Loop
			}

			notf := &notification.Notification{
				UserId:      d.acct.UserId,
				AccountId:   d.acct.Id,
				AccountType: digitalOcean,
				RemoteId:    strconv.Itoa(action.Id),
				Timestamp:   timestamp,
				Type:        action.Type,
				Subject:     subject,
			}

			if timestamp.Before(lastNotf.Timestamp) ||
				notf.RemoteId == lastNotf.RemoteId {
				break Loop
			}

			notfs = append(notfs, notf)
		}
	}

	pub := notification.NewPublisher(d.acct.UserId.Hex())
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

type DigitalOceanAuth struct{}

func (d *DigitalOceanAuth) Request(db *database.Database,
	userId bson.ObjectId) (url string, err error) {

	url, err = digitalOceanConf.Request(db, userId)
	if err != nil {
		return
	}

	return
}

func (d *DigitalOceanAuth) Authorize(db *database.Database, state string,
	code string) (acct *account.Account, err error) {

	coll := db.Accounts()

	auth, err := digitalOceanConf.Authorize(db, state, code)
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

func updateDigitalOcean() {
	digitalOceanConf = &oauth.Oauth2{
		Type:         digitalOcean,
		ClientId:     settings.DigitalOcean.ClientId,
		ClientSecret: settings.DigitalOcean.ClientSecret,
		CallbackUrl:  settings.System.Domain + "/callback/digitalocean",
		AuthUrl:      "https://cloud.digitalocean.com/v1/oauth/authorize",
		TokenUrl:     "https://cloud.digitalocean.com/v1/oauth/token",
		Scopes: []string{
			"read",
		},
	}
	digitalOceanConf.Config()
}
