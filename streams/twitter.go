package streams

import (
	"github.com/blckur/blckur/account"
	"github.com/blckur/blckur/database"
)

func (s *Stream) startTwitter() (err error) {
	db := database.GetDatabase()
	acctTyp, err := account.GetAccount(db, "", s.Id)
	if err != nil {
		return
	}
	acct := account.Twitter(*acctTyp)
	stream := acct.Stream(db)

	go func() {
		stream.Start()
	}()

	return
}
