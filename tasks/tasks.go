package tasks

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

type TaskManager interface {
	Push([]Task) (int, error)
	Pop(int) ([]Task, error)
}
