package settings

var System = &system{
	Id: "system",
}

func init() {
	register("system", System)
}

type system struct {
	Id        string `bson:"_id"`
	CookieKey []byte `bson:"cookie_key"`
	Domain    string `bson:"domain"`
	AppHome   string `bson:"app_home"`
}
