// Registered tasks with schedules.
package tasks

import (
	"github.com/blckur/blckur/database"
	"time"
)

var (
	registry [][][][]Task
)

type Task interface {
	Type() string
	Run(*database.Database) error
}

func register(hour int, min int, sec int, task Task) {
	if len(registry) == 0 {
		registry = make([][][][]Task, 25)

		for i := 0; i < 25; i++ {
			registry[i] = make([][][]Task, 61)
			for j := 0; j < 61; j++ {
				registry[i][j] = make([][]Task, 61)
				for k := 0; k < 61; k++ {
					registry[i][j][k] = []Task{}
				}
			}
		}
	}

	registry[hour+1][min+1][sec+1] = append(
		registry[hour+1][min+1][sec+1], task)
}

func registerMulti(hours []int, mins []int, secs []int, task Task) {
	for _, hour := range hours {
		for _, min := range mins {
			for _, sec := range secs {
				register(hour, min, sec, task)
			}
		}
	}
}

func GetTasks(t time.Time) (tasks []Task) {
	tasks = []Task{}

	for _, task := range registry[t.Hour()+1][t.Minute()+1][t.Second()+1] {
		tasks = append(tasks, task)
	}

	for _, task := range registry[0][t.Minute()+1][t.Second()+1] {
		tasks = append(tasks, task)
	}

	for _, task := range registry[0][0][t.Second()+1] {
		tasks = append(tasks, task)
	}

	for _, task := range registry[0][0][0] {
		tasks = append(tasks, task)
	}

	for _, task := range registry[t.Hour()+1][t.Minute()+1][0] {
		tasks = append(tasks, task)
	}

	for _, task := range registry[0][t.Minute()+1][0] {
		tasks = append(tasks, task)
	}

	for _, task := range registry[t.Hour()+1][0][0] {
		tasks = append(tasks, task)
	}

	return
}
