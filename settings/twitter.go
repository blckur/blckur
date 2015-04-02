package settings

var Twitter = &twitter{
	Id: "twitter",
}

func init() {
	register("twitter", Twitter)
}

type twitter struct {
	Id             string `bson:"_id"`
	ConsumerKey    string `bson:"consumer_key"`
	ConsumerSecret string `bson:"consumer_secret"`
}
