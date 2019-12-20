package pspk

type Keychain interface {
	WriteKey(name string, key []byte) (err error)
	ReadKey(name string) (key []byte, err error)
	Save() (err error)
}

type FileKeychain struct {
	basePath string
}

func (k *FileKeychain) WriteKey(name string, key []byte) (err error) {
	return
}

func (k *FileKeychain) ReadKey(name string) (key []byte, err error) {
	return
}

func (k *FileKeychain) Save() (err error) {
	return
}

func LoadKeychainFromDisk(basePath string) *FileKeychain {
	return &FileKeychain{
		basePath: basePath,
	}
}
