package account

import (
	"labix.org/v2/mgo/bson"
	"github.com/blckur/blckur/database"
	"labix.org/v2/mgo"
)

type Iter struct {
	coll *database.Collection
	iter *mgo.Iter
}

func (i *Iter) Iter() (out chan *Account) {
	out = make(chan *Account)

	go func() {
		acct := &Account{}
		for i.iter.Next(acct) {
			acct.ParseAlerts()

			out <- acct
			acct = &Account{}
		}

		close(out)
	}()

	return
}

func (i *Iter) Err() (err error) {
	err = i.iter.Err()
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func IterAccounts(db *database.Database, userId bson.ObjectId) (iter *Iter) {
	coll := db.Accounts()

	collIter := coll.Find(bson.M{
		"user_id": userId,
	}).Iter()

	iter = &Iter{
		coll: coll,
		iter: collIter,
	}

	return
}
