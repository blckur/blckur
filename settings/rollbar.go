package settings

var Rollbar = &rollbar{
	Id: "rollbar",
}

func init() {
	register("rollbar", Rollbar)
}

type rollbar struct {
	Id          string `bson:"_id"`
	Token       string `bson:"token"`
	Environment string `bson:"environment" default:"development"`
	RateLimit   int    `bson:"rate_limit" default:"300"`
}
