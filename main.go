package main

import (
	"github.com/blckur/blckur/cmd"
	"github.com/blckur/blckur/gdefer"
	"github.com/blckur/blckur/requires"
	"flag"
)

//	que := queue.NewQueue([]string{
//		"localhost:11100",
//		"localhost:11200",
//		"localhost:11300",
//		"localhost:11400",
//		"localhost:11500",
//		"localhost:11600",
//		"localhost:11700",
//		"localhost:11800",
//		"localhost:11900",
//	}, 3)
//
//	streams := que.GetStreams()
//
//	for _, stream := range streams {
//		go func(stream *queue.QueueStream) {
//			for {
//				job := stream.Reserve(6 * time.Hour)
//
//				print(".")
//
//				//fmt.Printf("%s: %s\n", job.Type, job.Data.(string))
//
//				err := job.Delete()
//				if err != nil {
//					panic(err)
//				}
//			}
//		}(stream)
//	}
//
//	db := database.GetDatabase()
//	defer db.Close()
//
//	err := messenger.Subscribe(db, "test", 1 * time.Second,
//			func(msg *messenger.Message) bool {
//		if msg == nil {
//			println("timeout")
//			return false
//		}
//
//		println("***************************************************")
//		println(msg.Data.(string))
//		println("***************************************************")
//		return false
//	})
//	if err != nil {
//		panic(err)
//	}

func main() {
	requires.Init()
	defer gdefer.End()
	flag.Parse()

	switch flag.Arg(0) {
	case "clear":
		cmd.Clear()
	case "set":
		cmd.Settings()
	case "app":
		cmd.App()
	case "queue":
		cmd.Queue()
	case "cache":
		cmd.Cache()
	case "scheduler":
		cmd.Scheduler()
	case "worker":
		cmd.Worker()
	}
}
