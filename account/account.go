package account

import (
	"labix.org/v2/mgo/bson"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/requires"
	"github.com/blckur/blckur/utils"
	"github.com/dropbox/godropbox/container/set"
	"labix.org/v2/mgo"
	"time"
)

var (
	eventTypes map[string]EventType = map[string]EventType{}
	acctTypes map[string][]string = map[string][]string{}
)

type Resource struct {
	Type string `bson:"type" json:"type"`
	Resource string `bson:"resource" json:"resource"`
	Events map[string]bool `bson:"events" json:"events"`
}

type EventType struct {
	Type string `bson:"-" json:"type"`
	ValueType string `bson:"-" json:"val_type"`
	Name string `bson:"-" json:"name"`
	Value interface{} `bson:"-" json:"value"`
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
	Events map[string]interface{} `bson:"events,omitempty" json:"-"`
	EventsParsed []*EventType `bson:"-" json:"events"`
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
	a.EventsParsed = []*EventType{}

	for _, evt := range acctTypes[a.Type] {
		typ, ok := eventTypes[evt]
		if !ok {
			continue
		}

		val := a.Events[evt]

		switch typ.ValueType{
		case "toggle":
			typ.Value = utils.InfToBool(val)
		}

		a.EventsParsed = append(a.EventsParsed, &typ)
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
	}).Iter()

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
