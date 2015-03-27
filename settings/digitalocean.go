package settings

var DigitalOcean = &digitalocean{
	Id: "digitalocean",
}

func init() {
	register("digitalocean", DigitalOcean)
}

type digitalocean struct {
	Id string `bson:"_id"`
	ClientId string `bson:"client_id"`
	ClientSecret string  `bson:"client_secret"`
	MaxMsg int `bson:"max_msg" default:"100"`
}
