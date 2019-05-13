package tasks

type TaskType int

const (
	UserTask TaskType = 0
	PostTask TaskType = 1
)

type Task struct {
	Type int
	Body string
	Deep int
}

type TaskManager interface {
	Push([]Task) (int, error)
	Pop(int) ([]Task, error)
}
