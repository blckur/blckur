// Oauth1 and oauth2 clients.
package oauth

import (
	"labix.org/v2/mgo/bson"
)

type Token struct {
	Id          string        `bson:"_id"`
	OauthSecret string        `bson:"oauth_secret"`
	Type        string        `bson:"type"`
	UserId      bson.ObjectId `bson:"user_id"`
}
