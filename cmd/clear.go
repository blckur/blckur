package cmd

import (
	"github.com/blckur/blckur/database"
	"labix.org/v2/mgo/bson"
)

// Clear cluster node info in database
func Clear() {
	db := database.GetDatabase()
	coll := db.Nodes()

	_, err := coll.RemoveAll(bson.M{})
	if err != nil {
		panic(err)
	}
}
