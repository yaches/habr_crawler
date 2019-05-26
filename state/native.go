package state

import (
	"github.com/yaches/habr_crawler/tasks"
)

type StorageNative struct {
	store map[tasks.Task]struct{}
}

func NewStorageNative() *StorageNative {
	return &StorageNative{store: make(map[tasks.Task]struct{})}
}

func (ss *StorageNative) Exists(task tasks.Task) (bool, error) {
	task.Deep = 0
	_, ok := ss.store[task]
	return ok, nil
}

func (ss *StorageNative) Add(task tasks.Task) error {
	task.Deep = 0
	ss.store[task] = struct{}{}
	return nil
}
