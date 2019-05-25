package state

type Storage interface {
	Exists(string) (bool, error)
	Add(string) error
}
