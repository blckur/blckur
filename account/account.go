package account

import (
	"labix.org/v2/mgo/bson"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/requires"
	"github.com/dropbox/godropbox/container/set"
	"time"
)

var (
	AlertTypes = map[string][]AlertType{
		"gmail": []AlertType{
			AlertType{
				Id: 1,
				Label: "All new messages",
				Type: "all",
			},
			AlertType{
				Id: 2,
				Label: "New messages matching sender",
				Type: "from",
				ValueType: "input",
				ValueLabel: "Enter complete or partial email address " +
				"of sender to match",
				ValueHolder: "Email address",
			},
			AlertType{
				Id: 3,
				Label: "New messages matching subject",
				Type: "subject",
				ValueType: "input",
				ValueLabel: "Enter search term to match in email subject",
				ValueHolder: "Search term",
			},
			AlertType{
				Id: 4,
				Label: "New messages matching message body",
				Type: "body",
				ValueType: "input",
				ValueLabel: "Enter search term to match in email body",
				ValueHolder: "Search term",
			},
		},
	}
	AlertLabels = map[string]map[string]string{
		"gmail": map[string]string{
			"all": "All messages",
			"from": "Messages matching sender",
			"subject": "Messages matching subject",
			"body": "Messages matching message body",
		},
	}
)

type Resource struct {
	Type string `bson:"type" json:"type"`
	Resource string `bson:"resource" json:"resource"`
	Events map[string]bool `bson:"events" json:"events"`
}

type AlertType struct {
	Id int `json:"id"`
	Label string `json:"label"`
	Type string `json:"type"`
	ValueType string `json:"value_type"`
	ValueLabel string `json:"value_label"`
	ValueHolder string `json:"value_holder"`
}

type Alert struct {
	Id string `bson:"_id" json:"id"`
	Type string `bson:"type" json:"type"`
	Value string `bson:"value" json:"value"`
	Label string `bson:"-" json:"label"`
}

type Account struct {
	Id bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserId bson.ObjectId `bson:"user_id,omitempty" json:"user_id"`
	Type string `bson:"type,omitempty" json:"type"`
	Identity string `bson:"identity,omitempty" json:"identity"`
	OauthTokn string `bson:"oauth_tokn,omitempty" json:"-"`
	OauthSec string `bson:"oauth_sec,omitempty" json:"-"`
	Oauth2AccTokn string `bson:"oauth2_acc_tokn,omitempty" json:"-"`
	Oauth2RefTokn string `bson:"oauth2_ref_tokn,omitempty" json:"-"`
	Oauth2Exp time.Time `bson:"oauth2_exp,omitempty" json:"-"`
	Alerts []*Alert `bson:"alerts,omitempty" json:"alerts"`
	Resources []*Resource `bson:"resources,omitempty" json:"resources"`
	coll *database.Collection `bson:"-" json:"-"`
}

func (a *Account) Commit() (err error) {
	err = a.coll.Commit(a.Id, a)
	return
}

func (a *Account) CommitFields(fields set.Set) (err error) {
	err = a.coll.CommitFields(a.Id, a, fields)
	return
}

func (a *Account) Marshal() {
	for _, alrt := range a.Alerts {
		if alrt.Id == nil {
			alrt.Id = bson.NewObjectId()
		}
		alrt.Label = AlertLabels[a.Type][alrt.Type]
	}
}

func (a *Account) Unmarshal() {
	a.Events = map[string]interface{}{}

	for _, evt := range a.EventsParsed {
		a.Events[evt.Type] = evt.Value
	}
}

func GetAccount(db *database.Database, userId bson.ObjectId,
		acctId bson.ObjectId) (acct *Account, err error) {
	coll := db.Accounts()
	acct = &Account{}

	err = coll.Find(bson.M{
		"_id": acctId,
		"user_id": userId,
	}).One(acct)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	acct.coll = coll

	return
}

func GetAccounts(db *database.Database, userId bson.ObjectId) (
		accts []*Account, err error) {
	coll := db.Accounts()
	accts = []*Account{}

	iter := coll.Find(bson.M{
		"user_id": userId,
	}).Sort("type", "identity").Iter()

	acct := &Account{}
	for iter.Next(acct) {
		acct.coll = coll
		acct.Marshal()

		accts = append(accts, acct)
		acct = &Account{}
	}

	err = iter.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func Init() {
	requires.After("settings")
	requires.Before("messenger")

	InitGmail()
	InitTwitter()

	requires.Register("account")
}
