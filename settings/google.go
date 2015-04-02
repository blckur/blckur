package settings

var Google = &google{
	Id: "google",
}

func init() {
	register("google", Google)
}

type google struct {
	Id           string `bson:"_id"`
	ClientId     string `bson:"client_id"`
	ClientSecret string `bson:"client_secret"`
	MaxMsg       int    `bson:"max_msg" default:"100"`
}
