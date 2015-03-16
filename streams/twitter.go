package streams

import (
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/settings"
	"time"
)

func (s *Stream) startTwitter() (err error) {
	db := database.GetDatabase()
	acctTyp, err := account.GetAccount(db, "", s.Id)
	if err != nil {
		return
	}
	acct := account.Twitter(*acctTyp)
	stream := acct.Stream(db)

	err = s.initialize()
	if err != nil {
		return
	}

	go func() {
		stream.Start()
	}()

	go func() {
		lastUpdate := time.Now()
		for {
			time.Sleep(time.Duration(
				settings.Stream.RefreshRate) * time.Second)
			stop, err := s.Update()
			if err != nil {
				panic(err)
			} else if stop {
				stream.Stop()
				return
			} else {
				lastUpdate = time.Now()
			}

			if time.Since(lastUpdate) > 2 * time.Minute {
				stream.Stop()
				return
			}
		}
	}()

	return
}
