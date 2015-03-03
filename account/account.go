package account

import (
	"labix.org/v2/mgo/bson"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/utils"
)

type Resource struct {
	Type string `bson:"type" json:"type"`
	Resource string `bson:"resource" json:"resource"`
	Events map[string]bool `bson:"events" json:"events"`
}

type Data struct {
	Id bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserId bson.ObjectId `bson:"user_id" json:"userid"`
	Type string `bson:"type" json:"type"`
	Name string `bson:"name" json:"name"`
	Token string `bson:"token" json:"-"`
	Secret string `bson:"secret" json:"-"`
	Events map[string]bool `bson:"events" json:"events"`
	Resources []*Resource `bson:"resources"`
}

type Account struct {
	*Data
	coll *database.Collection
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

	return
}

func Init() {
	utils.After("settings")
	utils.Before("messenger")

	InitTwitter()

	utils.Register("account")
}
