package pspk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/sah4ez/pspk/handler"
	"github.com/sah4ez/pspk/pkg/validation"
)

type PSPK interface {
	Publish(name string, key []byte) error
	Load(name string) ([]byte, error)
	GenerateLink(string) (string, error)
	DownloadByLink(string) (string, error)
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
	if err := validation.CheckLimitNameLen(name); err != nil {
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
	if err := validation.CheckLimitNameLen(name); err != nil {
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

func (p *pspk) GenerateLink(d string) (string, error) {
	if len(d) == 0 {
		return "", errors.New("empty data for generation link")
	}
	data := struct {
		Method string `json:"method"`
		Data   string `json:"data"`
	}{
		Method: handler.LinkKey,
		Data:   d,
	}

	b, err := json.Marshal(&data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", p.basePath, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	l := &struct {
		Link string `json:"link"`
	}{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&l)
	if err != nil {
		return "", err
	}
	return l.Link, nil
}

func (p *pspk) DownloadByLink(link string) (string, error) {
	req, err := http.NewRequest("GET", link, nil)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b := &struct {
		Data string `json:"data"`
	}{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&b)
	if err != nil {
		return "", err
	}
	return b.Data, nil
}
