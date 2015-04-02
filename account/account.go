// Interface to accounts collection in database.
package account

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/requires"
	"github.com/dropbox/godropbox/container/set"
	"github.com/dropbox/godropbox/errors"
	"labix.org/v2/mgo/bson"
	"reflect"
	"time"
)

type Client interface {
	Update(*database.Database) error
	Sync(*database.Database) error
	SetAccount(acct *Account)
}

type Auth interface {
	Request(*database.Database, bson.ObjectId) (string, error)
	Authorize(*database.Database, string, string) (*Account, error)
}

type Account struct {
	Id            bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserId        bson.ObjectId `bson:"user_id" json:"user_id"`
	Type          string        `bson:"type" json:"type"`
	Identity      string        `bson:"identity" json:"identity"`
	IdentityId    string        `bson:"identity_id,omitempty" json:"-"`
	New           bool          `bson:"new" json:"new"`
	OauthTokn     string        `bson:"oauth_tokn,omitempty" json:"-"`
	OauthSec      string        `bson:"oauth_sec,omitempty" json:"-"`
	Oauth2AccTokn string        `bson:"oauth2_acc_tokn,omitempty" json:"-"`
	Oauth2RefTokn string        `bson:"oauth2_ref_tokn,omitempty" json:"-"`
	Oauth2Exp     time.Time     `bson:"oauth2_exp,omitempty" json:"-"`
	Filters       []*Filter     `bson:"filters" json:"filters"`
}

func (a *Account) Commit(db *database.Database) (err error) {
	coll := db.Accounts()
	err = coll.Commit(a.Id, a)
	return
}

func (a *Account) CommitFields(db *database.Database,
	fields set.Set) (err error) {

	coll := db.Accounts()
	err = coll.CommitFields(a.Id, a, fields)

	return
}

func (a *Account) GetClient() (client Client, err error) {
	typ, ok := clientRegistry[a.Type]
	if !ok {
		err = &InvalidTypeError{
			errors.New("Invalid account type"),
		}
		return
	}

	val := reflect.New(typ).Elem()

	client = val.Addr().Interface().(Client)
	client.SetAccount(a)

	return
}

func GetAuth(acctType string) (auth Auth, authTyp int, err error) {
	typ, ok := authRegistry[acctType]
	if !ok {
		err = &InvalidTypeError{
			errors.New("Invalid account type"),
		}
		return
	}

	authTyp = authTypes[acctType]
	val := reflect.New(typ).Elem()

	auth = val.Addr().Interface().(Auth)

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
		acct.ParseFilters()

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
		for _, handler := range handlers {
			handler()
		}
	}
}
