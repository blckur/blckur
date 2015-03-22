package settings

var Stream = &stream{
	Id: "stream",
}

func init() {
	register("stream", Stream)
}

type stream struct {
	Id string `bson:"_id"`
	RefreshRate int `bson:"refresh_rate" default:"30"`
}
