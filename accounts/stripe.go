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
	"strings"
	"time"
)

const (
	stripe = "stripe"
)

var (
	stripeConf *oauth.Oauth2
)

func init() {
	account.Register(stripe, "Stripe", Oauth2,
		StripeAuth{}, StripeClient{},
		[]*account.FilterType{
			&account.FilterType{
				Label: "All new events",
				Type:  "all",
			},
			&account.FilterType{
				Label: "Charge succeeded",
				Type:  "charge_succeeded",
			},
			&account.FilterType{
				Label: "Charge failed",
				Type:  "charge_failed",
			},
			&account.FilterType{
				Label: "Charge refunded",
				Type:  "charge_refunded",
			},
			&account.FilterType{
				Label: "Chargeback created",
				Type:  "charge_dispute_created",
			},
			&account.FilterType{
				Label: "Chargeback closed",
				Type:  "charge_dispute_closed",
			},
			&account.FilterType{
				Label: "New customer",
				Type:  "customer_created",
			},
			&account.FilterType{
				Label: "Customer deleted",
				Type:  "customer_deleted",
			},
			&account.FilterType{
				Label: "Subscription created",
				Type:  "subscription_created",
			},
			&account.FilterType{
				Label: "Subscription deleted",
				Type:  "subscription_deleted",
			},
		}, func() {
			messenger.Register("settings", stripe,
				func(_ *messenger.Message) {
					updateStripe()
				})
			updateStripe()
		})
}

type stripeEvent struct {
	Id      string `json:"id"`
	Created int    `json:"created"`
	Type    string `json:"type"`
	Data    struct {
		Object struct {
			Id               string `json:"id"`
			Name             string `json:"name"`
			Email            string `json:"email"`
			Amount           int    `json:"amount"`
			Currency         string `json:"currency"`
			Reason           string `json:"reason"`
			Charge           string `json:"charge"`
			ChargeRefundable bool   `json:"is_charge_refundable"`
			Source           struct {
				Brand   string `json:"brand"`
				Funding string `json:"funding"`
				Country string `json:"country"`
			}
			Plan struct {
				Name string `json:"name"`
			} `json:"plan"`
		} `json:"object"`
	}
}

type StripeClient struct {
	acct *account.Account
}

func (s *StripeClient) SetAccount(acct *account.Account) {
	s.acct = acct
}

func (s *StripeClient) Update(db *database.Database) (err error) {
	client := stripeConf.NewClient(s.acct)

	data := struct {
		Email string `json:"email"`
	}{}

	err = client.GetJson("https://api.stripe.com/v1/account", &data)
	if err != nil {
		return
	}

	s.acct.Identity = data.Email

	return
}

func (s *StripeClient) filter(typ string) bool {
	for _, filter := range s.acct.Filters {
		if filter.Type == "all" || filter.Type == typ {
			return true
		}
	}

	return false
}

func (s *StripeClient) parse(evt *stripeEvent,
	lastNotf *notification.Notification, force bool) (
	notf *notification.Notification, stop bool) {

	stop = false
	timestamp := time.Unix(int64(evt.Created), 0)

	if lastNotf != nil && (evt.Id == lastNotf.RemoteId ||
		timestamp.Before(lastNotf.Timestamp)) {

		stop = true
		return
	}

	if force {
		notf = &notification.Notification{
			UserId:      s.acct.UserId,
			AccountId:   s.acct.Id,
			AccountType: stripe,
			RemoteId:    evt.Id,
			Timestamp:   timestamp,
		}
	}

	if lastNotf == nil {
		stop = true
		return
	}

	typ := strings.Replace(evt.Type, ".", "_", 1)

	var subject string
	var body string
	var link string

	switch evt.Type {
	case "charge.succeeded", "charge.failed", "charge.refunded":
		title := strings.Split(evt.Type, ".")[1]
		amount := float64(evt.Data.Object.Amount) / 100.
		subject = fmt.Sprintf("Charge %s for $%.2f", title, amount)
		link = fmt.Sprintf("https://dashboard.stripe.com/payments/%s",
			evt.Data.Object.Id)
		body = fmt.Sprintf("Charge %s from %s %s %s for $%.2f", title,
			evt.Data.Object.Source.Country, evt.Data.Object.Source.Brand,
			evt.Data.Object.Source.Funding, amount)

	case "charge.dispute.created", "charge.dispute.closed":
		title := strings.Split(evt.Type, ".")[2]
		amount := float64(evt.Data.Object.Amount) / 100.
		subject = fmt.Sprintf("Charge dispute %s for $%.2f", title, amount)
		link = fmt.Sprintf("https://dashboard.stripe.com/payments/%s",
			evt.Data.Object.Charge)

		var refund string
		if evt.Data.Object.ChargeRefundable {
			refund = "Refundable"
		} else {
			refund = "Non-refundable"
		}

		body = fmt.Sprintf("%s charge dispute %s for $%.2f reason %s",
			refund, title, amount, evt.Data.Object.Reason)

	case "customer.created", "customer.deleted":
		title := strings.Split(evt.Type, ".")[1]
		subject = fmt.Sprintf("Customer %s", title)
		body = fmt.Sprintf("Customer %s with email %s", title,
			evt.Data.Object.Email)
		link = fmt.Sprintf("https://dashboard.stripe.com/customers/%s",
			evt.Data.Object.Id)

	case "subscription_created", "subscription_deleted":
		title := strings.Split(evt.Type, ".")[1]
		subject = fmt.Sprintf("Subscription %s", title)
		body = fmt.Sprintf("Subscription %s has been %s",
			evt.Data.Object.Email, title)
		link = fmt.Sprintf("https://dashboard.stripe.com/subscriptions/%s",
			evt.Data.Object.Id)

	default:
		return
	}

	if !s.filter(typ) {
		return
	}

	notf = &notification.Notification{
		UserId:      s.acct.UserId,
		AccountId:   s.acct.Id,
		AccountType: stripe,
		RemoteId:    evt.Id,
		Timestamp:   timestamp,
		Type:        typ,
		Link:        link,
		Subject:     subject,
		Body:        body,
	}

	return
}

func (s *StripeClient) Sync(db *database.Database) (err error) {
	client := stripeConf.NewClient(s.acct)

	err = client.Refresh(db)
	if err != nil {
		return
	}

	lastNotf, err := notification.GetLastNotification(db,
		s.acct.UserId, s.acct.Id)
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

	n := settings.Stripe.MaxMsg / 10
Loop:
	for i := 0; i < n; i++ {
		events := struct {
			Data []*stripeEvent `json:"data"`
		}{}

		url := fmt.Sprintf("https://api.stripe.com/v1/events?limit=%d",
			msgCount)

		if pageToken != "" {
			url += "&starting_after=" + pageToken
		}

		err = client.GetJson(url, &events)
		if err != nil {
			return
		}

		for j, evt := range events.Data {
			if j == 0 {
				pageToken = evt.Id
			}

			notf, stop := s.parse(evt, lastNotf, i == 0 && j == 0)
			if notf != nil {
				notfs = append(notfs, notf)
			}

			if stop {
				break Loop
			}
		}
	}

	pub := notification.NewPublisher(s.acct.UserId.Hex())
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

type StripeAuth struct{}

func (s *StripeAuth) Request(db *database.Database,
	userId bson.ObjectId) (url string, err error) {

	url, err = stripeConf.Request(db, userId)
	if err != nil {
		return
	}

	return
}

func (s *StripeAuth) Authorize(db *database.Database, state string,
	code string) (acct *account.Account, err error) {

	coll := db.Accounts()

	auth, err := stripeConf.Authorize(db, state, code)
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

func updateStripe() {
	stripeConf = &oauth.Oauth2{
		Type:         stripe,
		ClientId:     settings.Stripe.ClientId,
		ClientSecret: settings.Stripe.ClientSecret,
		CallbackUrl:  settings.System.Domain + "/callback/stripe",
		AuthUrl:      "https://connect.stripe.com/oauth/authorize",
		TokenUrl:     "https://connect.stripe.com/oauth/token",
		Scopes: []string{
			"read_only",
		},
	}
	stripeConf.Config()
}
