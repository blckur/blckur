package account

import (
	"labix.org/v2/mgo/bson"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/requires"
	"time"
)

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
	OauthTokn string `bson:"oauth_tokn,omitempty" json:"-"`
	OauthSec string `bson:"oauth_sec,omitempty" json:"-"`
	Oauth2AccTokn string `bson:"oauth2_acc_tokn,omitempty" json:"-"`
	Oauth2RefTokn string `bson:"oauth2_ref_tokn,omitempty" json:"-"`
	Oauth2Exp time.Time `bson:"oauth2_exp,omitempty" json:"-"`
	Events map[string]bool `bson:"events" json:"events"`
	Resources []*Resource `bson:"resources" json:"resources"`
	coll *database.Collection `bson:"-" json:"-"`
}

type Token struct {
	Token string `bson:"_id"`
	Secret string `bson:"secret"`
	UserId bson.ObjectId `bson:"user_id"`
}

func GetAccounts(db *database.Database, userId bson.ObjectId) (
		accts []*Account, err error) {
	coll := db.Accounts()
	accts = []*Account{}

	coll.Find(bson.M{
		"user_id": userId,
	}).All(&accts)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	for _, acct := range accts {
		acct.coll = coll
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
