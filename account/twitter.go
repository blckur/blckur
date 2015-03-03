package account

import (
	"encoding/json"
	"github.com/blckur/blckur/database"
	"github.com/blckur/blckur/errortypes"
	"github.com/blckur/blckur/settings"
	"github.com/blckur/blckur/utils"
	"github.com/blckur/blckur/messenger"
	"github.com/blckur/blckur/logger"
	"github.com/blckur/blckur/constants"
	"github.com/dropbox/godropbox/errors"
	"github.com/mrjones/oauth"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"time"
)

var (
	consumer *oauth.Consumer
	callback string
)

type Twitter struct {
	*Account
}

func AuthTwitter(db *database.Database, userId bson.ObjectId) (
		url string, err error) {
	coll := db.Tokens()

	reqToken, url, err := consumer.GetRequestTokenAndUrl(callback)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "account: Unknown twitter api error"),
		}
		return
	}

	tokn := &Token{
		Token: reqToken.Token,
		Secret: reqToken.Secret,
		UserId: userId,
	}

	err = coll.Insert(tokn)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func NewTwitter(db *database.Database, token string, verifier string) (
	acct *Account, err error) {
	toknColl := db.Tokens()
	acctColl := db.Accounts()
	tokn := &Token{}

	err = toknColl.FindOneId(token, tokn)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	reqTokn := &oauth.RequestToken{
		Token: tokn.Token,
		Secret: tokn.Secret,
	}

	accessTokn, err := consumer.AuthorizeToken(reqTokn, verifier)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "account: Unknown twitter api error"),
		}
		return
	}

	resp, err := consumer.Get(
		"https://api.twitter.com/1.1/account/settings.json",
		nil, accessTokn)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "account: Unknown twitter api error"),
		}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "account: Unknown parse error"),
		}
		return
	}

	data := map[string]interface{}{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "account: Unknown parse error"),
		}
		return
	}

	name := data["screen_name"].(string)

	acct = &Account{
		&Data{
			UserId: tokn.UserId,
			Type: "twitter",
			Name: name,
			Token: tokn.Token,
			Secret: tokn.Secret,
		},
		acctColl,
	}

	err = acctColl.Insert(acct.Data)
	if err != nil {
		err = database.ParseError(err)
		return
	}

	return
}

func update() {
	for {
		err := settings.Twitter.Update()
		if err != nil {
			logger.Error("account: Twitter update error %s", err)
		} else {
			break
		}

		time.Sleep(constants.DB_RETRY_DELAY)
	}

	consumer = oauth.NewConsumer(
		settings.Twitter.ConsumerKey,
		settings.Twitter.ConsumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl: "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl: "https://api.twitter.com/oauth/access_token",
		},
	)

	callback = settings.System.Domain + "/callback/twitter"
}

func Init() {
	utils.After("settings")
	utils.Before("messenger")

	update()
	messenger.Register("settings", "twitter", update)

	utils.Register("account")
}
