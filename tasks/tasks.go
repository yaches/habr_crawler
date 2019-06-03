package tasks

import (
	"fmt"
)

type TaskType int

const (
	PostTask      TaskType = 0
	UserTask      TaskType = 1
	UserPostsTask TaskType = 2
)

type Task struct {
	Type TaskType
	// if type is UserTask or UserPostsTask, body is username; if type is PostTask, body is post id.
	Body string
	Deep int
	// Page is page number for UserPostsTask
	Page int
}

type Manager interface {
	Channel() <-chan Task
	Push(Task) error
	Done(Task) error
}

const fstr = "%d;%d;%s"

func Encode(task Task) (string, error) {
	return fmt.Sprintf(fstr, task.Type, task.Page, task.Body), nil
}

func Decode(str string) (Task, error) {
	task := Task{}
	_, err := fmt.Sscanf(str, fstr, &task.Type, &task.Page, &task.Body)
	return task, err
}
