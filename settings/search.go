package settings

var Search = &search{
	Id: "search",
}

func init() {
	register("search", Search)
}

type search struct {
	Id      string `bson:"_id"`
	Address string `bson:"address"`
}
