package settings

var Account = &account{
	Id: "account",
}

func init() {
	register("account", Account)
}

type account struct {
	Id         string `bson:"_id"`
	UpdateRate int    `bson:"update_rate" default:"60"`
	BodyLimit  int    `bson:"body_limit" default:"300"`
}
