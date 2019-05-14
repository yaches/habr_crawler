package statestore

type StorageNative struct {
	store map[string]struct{}
}

func NewStorageNative() StorageNative {
	return StorageNative{store: make(map[string]struct{})}
}

func (ss *StorageNative) Exists(taskBody string) (bool, error) {
	_, ok := ss.store[taskBody]
	return ok, nil
}

func (ss *StorageNative) Add(taskBody string) error {
	ss.store[taskBody] = struct{}{}
	return nil
}
