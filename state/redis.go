package state

import (
	"github.com/go-redis/redis"
	"github.com/yaches/habr_crawler/tasks"
)

type StorageRedis struct {
	db *redis.Client
}

func NewStorageRedis(db *redis.Client) *StorageRedis {
	return &StorageRedis{db: db}
}

func (s *StorageRedis) Add(task tasks.Task) error {
	task.Deep = 0
	str, err := tasks.Encode(task)
	if err != nil {
		return err
	}
	return s.db.Set(str, "", 0).Err()
}

func (s *StorageRedis) Exists(task tasks.Task) (bool, error) {
	task.Deep = 0
	str, err := tasks.Encode(task)
	if err != nil {
		return true, err
	}
	r := s.db.Exists(str)
	return r.Val() != 0, r.Err()
}

func (s *StorageRedis) Close() {
	s.db.Close()
}
