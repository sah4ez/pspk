package utils

import (
	"io/ioutil"
	"os"
)

func Write(path, name string, data []byte) error {
	os.Mkdir(path, 0766)
	f, err := os.OpenFile(path+"/"+name, os.O_CREATE|os.O_RDWR, 0766)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

func Read(path, name string) ([]byte, error) {
	return ReadPath(path + "/" + name)
}

func ReadPath(path string) ([]byte, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0766)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return b, nil
}
