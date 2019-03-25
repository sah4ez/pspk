package pspk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"unicode/utf8"
)

type PSPK interface {
	Publish(name string, key []byte) error
	Load(name string) ([]byte, error)
	Link([]byte) (string, error)
}

type pspk struct {
	client   *http.Client
	basePath string
}

func New(basePath string) *pspk {
	return &pspk{
		client:   &http.Client{},
		basePath: basePath,
	}
}

func (p *pspk) Publish(name string, key []byte) error {
	if err := checkLimitNameLen(name); err != nil {
		return err
	}
	body := struct {
		Name string `json:"name"`
		Key  string `json:"key"`
	}{
		Name: name,
		Key:  base64.StdEncoding.EncodeToString(key),
	}

	b, err := json.Marshal(&body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", p.basePath, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (p *pspk) Load(name string) ([]byte, error) {
	if err := checkLimitNameLen(name); err != nil {
		return nil, err
	}
	body := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}

	b, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", p.basePath, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	key := &struct {
		Key string `json:"key"`
	}{}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &key)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(key.Key)
}

func (p *pspk) Link([]byte) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func checkLimitNameLen(name string) error {
	if utf8.RuneCountInString(name) > 1000 {
		return fmt.Errorf("limit up to 1000 sign for name of key")
	}
	return nil
}
