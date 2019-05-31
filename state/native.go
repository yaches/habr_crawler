package state

import (
	"sync"

	"github.com/yaches/habr_crawler/tasks"
)

type StorageNative struct {
	store map[tasks.Task]struct{}
	mx    sync.Mutex
}

func NewStorageNative() *StorageNative {
	return &StorageNative{store: make(map[tasks.Task]struct{}), mx: sync.Mutex{}}
}

func (ss *StorageNative) Exists(task tasks.Task) (bool, error) {
	ss.mx.Lock()
	defer ss.mx.Unlock()
	task.Deep = 0
	_, ok := ss.store[task]
	return ok, nil
}

func (ss *StorageNative) Add(task tasks.Task) error {
	ss.mx.Lock()
	defer ss.mx.Unlock()
	task.Deep = 0
	ss.store[task] = struct{}{}
	return nil
}
