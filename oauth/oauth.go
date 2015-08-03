// Oauth1 and oauth2 clients.
package oauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/blckur/blckur/utils"
	"labix.org/v2/mgo/bson"
	"math/rand"
	"strconv"
	"time"
)

type Token struct {
	Id          string        `bson:"_id"`
	OauthSecret string        `bson:"oauth_secret"`
	Type        string        `bson:"type"`
	UserId      bson.ObjectId `bson:"user_id"`
}

func (c *Oauth1Client) Sign(method string, url string,
	params map[string]string) (sig string) {
	message := method + "&" + utils.Escape(url)
	delimEsc := utils.Escape("&")
	first := true

	addParam := func(key string, val string) {
		if first {
			first = false
			message += "&"
		} else {
			message += delimEsc
		}
		message += utils.Escape(fmt.Sprintf("%s=%s", key, val))
	}

	addParam("oauth_consumer_key", c.conf.ConsumerKey)
	addParam("oauth_nonce", strconv.FormatInt(rand.Int63(), 10))
	addParam("oauth_signature_method", "HMAC-SHA1")
	addParam("oauth_timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	addParam("oauth_token", c.Token)
	addParam("oauth_version", "1.0")

	for key, val := range params {
		addParam(key, val)
	}

	key := utils.Escape(c.conf.ConsumerSecret) + "&" + utils.Escape(c.Secret)

	hashFunc := hmac.New(sha1.New, []byte(key))
	hashFunc.Write([]byte(message))
	rawSignature := hashFunc.Sum(nil)
	sig = base64.StdEncoding.EncodeToString(rawSignature)

	return
}
