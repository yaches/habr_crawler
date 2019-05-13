package statestore

type StateStoreNative struct {
	store map[string]struct{}
}

func NewStateStoreNative() StateStoreNative {
	return StateStoreNative{store: make(map[string]struct{})}
}

func (ss *StateStoreNative) Exists(taskBody string) (bool, error) {
	_, ok := ss.store[taskBody]
	return ok, nil
}

func (ss *StateStoreNative) Add(taskBody string) error {
	ss.store[taskBody] = struct{}{}
	return nil
}
