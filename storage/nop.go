package storage

type nopStorage struct{}

func (n nopStorage) LoadFile(_ string) ([]byte, error) {
	return nil, nil
}

func (n nopStorage) ListFiles() ([]string, error) {
	return nil, nil
}

func (n nopStorage) StoreFile(_ string, _ []byte) error {
	return nil
}

func newNopStorage() Storage {
	return nopStorage{}
}
