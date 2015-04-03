package settings

var GitHub = &github{
	Id: "github",
}

func init() {
	register("github", GitHub)
}

type github struct {
	Id           string `bson:"_id"`
	ClientId     string `bson:"client_id"`
	ClientSecret string `bson:"client_secret"`
}
