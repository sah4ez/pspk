// +build js,wasm

package utils

import (
	"fmt"
	"sync"
)

type WasmStorage struct {
	m      sync.RWMutex
	bucket map[string][]byte
}

func NewWasmStorage() *WasmStorage {
	return &WasmStorage{
		bucket: map[string][]byte{},
	}
}

func (fs *WasmStorage) Write(path, name string, data []byte) error {
	// fs.m.Lock()
	// defer fs.m.Unlock()

	fmt.Println("write", path+"/"+name)
	fs.bucket[path+"/"+name] = data
	return nil
}

func (fs *WasmStorage) Read(path, name string) ([]byte, error) {
	// fs.m.RLock()
	// defer fs.m.RUnlock()

	fmt.Println("read", path+"/"+name)
	if d, ok := fs.bucket[path+"/"+name]; ok {
		return d, nil
	}
	return nil, fmt.Errorf("not found key: %s", path+"/"+name)
}
