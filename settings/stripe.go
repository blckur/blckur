package settings

var Stripe = &stripe{
	Id: "stripe",
}

func init() {
	register("stripe", Stripe)
}

type stripe struct {
	Id           string `bson:"_id"`
	ClientId     string `bson:"client_id"`
	ClientSecret string `bson:"client_secret"`
	MaxMsg       int    `bson:"max_msg" default:"100"`
}
