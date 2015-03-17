package cmd

import (
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/logger"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/utils"
	"github.com/blckur/blckur/messenger"
	"flag"
)

func Settings() {
	utils.SeedRand()

	logger.Init()
	database.Init()
	settings.Init()

	group := flag.Arg(1)
	key := flag.Arg(2)
	val := flag.Arg(3)
	db := database.GetDatabase()
	defer db.Close()

	err := settings.Set(db, group, key, val)
	if err != nil {
		panic(err)
	}

	err = messenger.Publish(db, "settings", group)
	if err != nil {
		panic(err)
	}

	return
}