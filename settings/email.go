package settings

var Email = &email{
	Id: "email",
}

func init() {
	register("email", Email)
}

type email struct {
	Id       string `bson:"_id"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	Server   string `bson:"server"`
	Port     int    `bson:"port"`
	From     string `bson:"from"`
}
