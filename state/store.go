package state

import (
	"github.com/yaches/habr_crawler/tasks"
)

type Storage interface {
	Exists(tasks.Task) (bool, error)
	Add(tasks.Task) error
}
