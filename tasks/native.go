package tasks

import "errors"

type TaskManagerChan struct {
	ch chan Task
}

func NewTaskManagerChan() *TaskManagerChan {
	return &TaskManagerChan{ch: make(chan Task, 1000)}
}

func (tm *TaskManagerChan) Push(tasks []Task) (int, error) {
	if tasks == nil {
		return 0, errors.New("Tasks is nil")
	}

	i := 0
	for _, t := range tasks {
		tm.ch <- t
		i++
	}

	return i, nil
}

func (tm *TaskManagerChan) Pop(n int) ([]Task, error) {
	var tasks []Task
	var t Task
	for i := 0; i < n; i++ {
		t = <-tm.ch
		tasks = append(tasks, t)
	}

	return tasks, nil
}
