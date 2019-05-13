package statestore

type StateStore interface {
	Exists(string) (bool, error)
	Add(string) error
}
