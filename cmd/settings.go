package cmd

import (
	"flag"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/settings"
	"strconv"
)

// Manages system settigs stored in database
func Settings() {
	group := flag.Arg(1)
	key := flag.Arg(2)
	val := flag.Arg(3)
	db := database.GetDatabase()
	defer db.Close()

	var valParsed interface{}
	if x, err := strconv.Atoi(val); err == nil {
		valParsed = x
	} else {
		valParsed = val
	}

	err := settings.Set(db, group, key, valParsed)
	if err != nil {
		panic(err)
	}

	err = messenger.Publish(db, "settings", group)
	if err != nil {
		panic(err)
	}

	return
}
