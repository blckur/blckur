package settings

var Queue = &queue{
	Id: "queue",
}

func init() {
	register("queue", Queue)
}

type queue struct {
	Id string `bson:"_id"`
	Consistency int `bson:"consistency" default:"2"`
}
