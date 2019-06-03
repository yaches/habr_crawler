package tasks

type ManagerNative struct {
	ch chan Task
}

func NewManagerNative() *ManagerNative {
	return &ManagerNative{ch: make(chan Task, 1000)}
}

func (m *ManagerNative) Channel() <-chan Task {
	return m.ch
}

func (m *ManagerNative) Push(task Task) error {
	m.ch <- task
	return nil
}

func (m *ManagerNative) Done(task Task) error {
	return nil
}
