package settings

var PapperTrail = &papperTrail{
	Id: "paper_trail",
}

func init() {
	register("paper_trail", PapperTrail)
}

type papperTrail struct {
	Id        string `bson:"_id"`
	Address   string `bson:"address"`
	RateLimit int    `bson:"rate_limit" default:"180"`
}
