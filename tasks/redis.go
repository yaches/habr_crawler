package tasks

import (
	"github.com/go-redis/redis"
)

const (
	channel = "habr_tasks"
	bufSz   = 1000000000000
)

type ManagerRedis struct {
	t_ch chan Task
	r_ch <-chan *redis.Message
	db   *redis.Client
}

func NewManagerRedis(db *redis.Client) *ManagerRedis {
	m := &ManagerRedis{
		db:   db,
		r_ch: db.Subscribe(channel).ChannelSize(bufSz),
		t_ch: make(chan Task, bufSz),
	}

	go func() {
		for {
			mes := <-m.r_ch
			task, err := Decode(mes.Payload)
			if err != nil {
				continue
			}
			m.t_ch <- task
		}
	}()

	return m
}

func (m *ManagerRedis) Channel() <-chan Task {
	return m.t_ch
}

func (m *ManagerRedis) Push(task Task) error {
	str, err := Encode(task)
	if err != nil {
		return err
	}
	return m.db.Publish(channel, str).Err()
}
