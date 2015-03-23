package settings

var Task = &task{
	Id: "task",
}

func init() {
	register("task", Task)
}

type task struct {
	Id string `bson:"_id"`
	Timeout int `bson:"timeout" default:"45"`
	RetryTimeout int `bson:"timeout" default:"15"`
}
