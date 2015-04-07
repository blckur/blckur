package accounts

import (
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/oauth"
	"github.com/blckur/blckur/settings"
	"labix.org/v2/mgo/bson"
)

var (
	stripeConf *oauth.Oauth2
)

func init() {
	account.Register("stripe", "Stripe", OAUTH2,
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
			messenger.Register("settings", "stripe",
				func(_ *messenger.Message) {
					updateStripe()
				})
			updateStripe()
		})
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
		Login string `json:"login"`
	}{}

	err = client.GetJson("", &data)
	if err != nil {
		return
	}

	s.acct.Identity = data.Login

	return
}

func (s *StripeClient) Sync(db *database.Database) (err error) {
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
		Type:         "stripe",
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
