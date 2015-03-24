package account

import (
	"labix.org/v2/mgo/bson"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/requires"
	"github.com/dropbox/godropbox/container/set"
	"time"
)

type clientInterface interface {
	Update(*database.Database) error
	Sync(*database.Database) error
}

type Resource struct {
	Type string `bson:"type" json:"type"`
	Resource string `bson:"resource" json:"resource"`
	Events map[string]bool `bson:"events" json:"events"`
}

type Account struct {
	Id bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserId bson.ObjectId `bson:"user_id" json:"user_id"`
	Type string `bson:"type" json:"type"`
	Identity string `bson:"identity" json:"identity"`
	IdentityId string `bson:"identity_id,omitempty" json:"-"`
	OauthTokn string `bson:"oauth_tokn,omitempty" json:"-"`
	OauthSec string `bson:"oauth_sec,omitempty" json:"-"`
	Oauth2AccTokn string `bson:"oauth2_acc_tokn,omitempty" json:"-"`
	Oauth2RefTokn string `bson:"oauth2_ref_tokn,omitempty" json:"-"`
	Oauth2Exp time.Time `bson:"oauth2_exp,omitempty" json:"-"`
	Alerts []*Alert `bson:"alerts" json:"alerts"`
	Resources []*Resource `bson:"resources" json:"resources"`
	coll *database.Collection
}

func (a *Account) Commit() (err error) {
	err = a.coll.Commit(a.Id, a)
	return
}

func (a *Account) CommitFields(fields set.Set) (err error) {
	err = a.coll.CommitFields(a.Id, a, fields)
	return
}

func (a *Account) Update(db *database.Database) (err error) {
	client := a.getClientInterface()

	err = client.Update(db)
	if err != nil {
		return
	}

	return
}

func (a *Account) Sync(db *database.Database) (err error) {
	client := a.getClientInterface()

	err = client.Sync(db)
	if err != nil {
		return
	}

	return
}

func (a *Account) getClientInterface() (clientIntf clientInterface) {
	if a.Type == "twitter" {
		a := twitter(*a)
		clientIntf = &a
	} else if a.Type == "gmail" {
		a := gmail(*a)
		clientIntf = &a
	}

	return
}

func GetAccount(db *database.Database, userId bson.ObjectId,
		acctId bson.ObjectId) (acct *Account, err error) {
	coll := db.Accounts()
	acct = &Account{}

	query := bson.M{
		"_id": acctId,
	}

	if userId != "" {
		query["user_id"] = userId
	}

	err = coll.Find(query).One(acct)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	acct.coll = coll

	return
}

func RemAccount(db *database.Database, userId bson.ObjectId,
		acctId bson.ObjectId) (err error) {
	coll := db.Accounts()

	query := bson.M{
		"_id": acctId,
	}

	if userId != "" {
		query["user_id"] = userId
	}

	err = coll.Remove(query)
	if err != nil {
		err = database.ParseError(err)
		return
	}

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
		acct.ParseAlerts()

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

func init() {
	module := requires.New("account")
	module.After("settings")
	module.Before("messenger")

	module.Handler = func() {
		initGmail()
		initTwitter()
	}
}
