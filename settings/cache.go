package settings

var Cache = &cache{
	Id: "cache",
}

func init() {
	register("cache", Cache)
}

type cache struct {
	Id string `bson:"_id"`
	Consistency int `bson:"consistency" default:"2"`
	TimeoutMilli int `bson:"timeout_milli" default:"2500"`
	MaxIdle int `bson:"max_idle" default:"3"`
	MaxActive int `bson:"max_active" default:"0"`
	IdleTimeout int `bson:"idle_timeout" default:"300"`
}
