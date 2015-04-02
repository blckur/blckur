// Registered tasks with schedules.
package tasks

import (
	"github.com/blckur/blckur/database"
	"time"
)

var (
	registry [][][]Task
)

type Task interface {
	Type() string
	Run(*database.Database) error
}

func register(hour int, min int, task Task) {
	if len(registry) == 0 {
		registry = make([][][]Task, 25)

		for i := 0; i < 25; i++ {
			registry[i] = make([][]Task, 61)
			for j := 0; j < 61; j++ {
				registry[i][j] = []Task{}
			}
		}
	}

	registry[hour+1][min+1] = append(
		registry[hour+1][min+1],
		task)
}

func registerMulti(hours []int, mins []int, task Task) {
	for hour := range hours {
		for min := range mins {
			register(hour, min, task)
		}
	}
}

func GetTasks(tme time.Time) (tasks []Task) {
	tasks = []Task{}

	for _, task := range registry[tme.Hour()+1][tme.Minute()+1] {
		tasks = append(tasks, task)
	}

	for _, task := range registry[0][tme.Minute()+1] {
		tasks = append(tasks, task)
	}

	for _, task := range registry[tme.Hour()+1][0] {
		tasks = append(tasks, task)
	}

	for _, task := range registry[0][0] {
		tasks = append(tasks, task)
	}

	return
}
