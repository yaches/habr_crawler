package tasks

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

const (
	channel = "habr_tasks"
	bufSz   = 100000000
)

type ManagerRedis struct {
	ch chan Task
	db *redis.Client
}

func NewManagerRedis(db *redis.Client) *ManagerRedis {
	m := &ManagerRedis{
		db: db,
		ch: make(chan Task, bufSz),
	}
	return m
}

func (m *ManagerRedis) Channel() <-chan Task {
	return m.ch
}

func (m *ManagerRedis) Push(task Task) error {
	str, err := Encode(task)
	if err != nil {
		return err
	}

	exRes, err := m.db.Exists(str).Result()
	if err != nil {
		return err
	}

	if exRes == 0 {
		err = m.db.Set(str, "0", 0).Err()
		if err != nil {
			return err
		}
		m.ch <- task
	}

	return nil
}

func (m *ManagerRedis) Done(task Task) error {
	str, err := Encode(task)
	if err != nil {
		return err
	}

	return m.db.Set(str, "1", 0).Err()
}

func (m *ManagerRedis) Fill() error {
	var cur uint64
	for {
		r := m.db.Scan(cur, "*", 1000)
		if r.Err() != nil {
			return r.Err()
		}
		for i := r.Iterator(); i.Next(); {
			v, err := m.db.Get(i.Val()).Result()
			if err != nil || i.Err() != nil {
				zap.L().Warn("Redis:", zap.Error(err), zap.Error(i.Err()))
				continue
			}
			if v == "0" {
				task, err := Decode(i.Val())
				if err != nil {
					zap.L().Warn("", zap.Error(err))
					continue
				}
				m.ch <- task
			}
		}
		if cur == 0 {
			break
		}
	}
	return nil
}
