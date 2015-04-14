package settings

var Account = &account{
	Id: "account",
}

func init() {
	register("account", Account)
}

type account struct {
	Id        string `bson:"_id"`
	BodyLimit int    `bson:"body_limit" default:"300"`
}
