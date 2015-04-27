package settings

var Rollbar = &rollbar{
	Id: "rollbar",
}

func init() {
	register("rollbar", Rollbar)
}

type rollbar struct {
	Id    string `bson:"_id"`
	Token string `bson:"token"`
}
