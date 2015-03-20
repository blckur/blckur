package tasks

import (
	"github.com/blckur/blckur/database"
	"time"
)

var (
	tasksRegistry [][][]Task
)

type Task interface {
	Type() string
	Run(db *database.Database)
}

func register(hour int, min int, task Task) {
	if len(tasksRegistry) == 0 {
		tasksRegistry = make([][][]Task, 25)

		for i := 0; i < 25; i++ {
			tasksRegistry[i] = make([][]Task, 61)
			for j := 0; j < 61; j++ {
				tasksRegistry[i][j] = []Task{}
			}
		}
	}

	tasksRegistry[hour + 1][min + 1] = append(
		tasksRegistry[hour + 1][min + 1],
		task)
}

func registerMulti(hours []int, mins int, task Task) {
	for hour := range hours {
		for min := range mins {
			register(hour, min, task)
		}
	}
}

func GetTasks(tme time.Time) (tasks []Task) {
	tasks = []Task{}

	for _, task := range tasksRegistry[tme.Hour() + 1][tme.Minute() + 1] {
		tasks = append(tasks, task)
	}

	for _, task := range tasksRegistry[0][tme.Minute() + 1] {
		tasks = append(tasks, task)
	}

	for _, task := range tasksRegistry[tme.Hour() + 1][0] {
		tasks = append(tasks, task)
	}

	for _, task := range tasksRegistry[0][0] {
		tasks = append(tasks, task)
	}

	return
}
