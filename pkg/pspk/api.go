package pspk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/sah4ez/pspk/pkg/validation"
)

const (
	UserAgent     = "pspk-client/1.0.0"
	NameKey       = "name_key"
	NameSearchKey = "name_regex"
	LinkKey       = "link"
	OutputKey     = "output"
	LastIDKEy     = "last_key"
	LimitKey      = "limit"
	QRCodeKey     = "qr_code"
)

type Key struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type GetAllOptions struct {
	Output  string `url:"output,omitempty"`
	LastKey string `url:"last_key,omitempty"`
	Limit   int    `url:"limit,omitempty"`
}

type PSPK interface {
	Publish(name string, key []byte) error
	Load(name string) ([]byte, error)
	GenerateLink(string) (string, error)
	DownloadByLink(string) (string, error)
	GetAll(opts GetAllOptions) ([]Key, error)
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
	req, err := p.newRequest("POST", body, nil)
	if err != nil {
		return err
	}
	return p.doRequest(req, nil)
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
	req, err := p.newRequest("POST", body, nil)
	if err != nil {
		return nil, err
	}
	resource := struct {
		Key string `json:"key"`
	}{}
	err = p.doRequest(req, &resource)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(resource.Key)
}

func (p *pspk) GenerateLink(d string) (string, error) {
	if len(d) == 0 {
		return "", errors.New("empty data for generation link")
	}
	body := struct {
		Method string `json:"method"`
		Data   string `json:"data"`
	}{
		Method: LinkKey,
		Data:   d,
	}
	req, err := p.newRequest("POST", body, nil)
	if err != nil {
		return "", err
	}
	resource := struct {
		Link string `json:"link"`
	}{}
	err = p.doRequest(req, &resource)
	if err != nil {
		return "", err
	}
	return resource.Link, nil
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

func (p *pspk) GetAll(options GetAllOptions) ([]Key, error) {
	req, err := p.newRequest("GET", nil, options)
	if err != nil {
		return nil, err
	}
	var resource []Key
	err = p.doRequest(req, &resource)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

func (p *pspk) newRequest(method string, body, options interface{}) (*http.Request, error) {
	u, err := url.Parse(p.basePath)
	if err != nil {
		return nil, err
	}
	// Custom options
	if options != nil {
		optionsQuery, err := query.Values(options)
		if err != nil {
			return nil, err
		}
		for k, values := range u.Query() {
			for _, v := range values {
				optionsQuery.Add(k, v)
			}
		}
		u.RawQuery = optionsQuery.Encode()
	}
	var js []byte = nil
	if body != nil {
		js, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(js))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", UserAgent)
	return req, nil
}

func (p *pspk) doRequest(req *http.Request, v interface{}) error {
	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = checkResponseError(resp)
	if err != nil {
		return err
	}
	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return err
		}
	}
	return nil
}

func checkResponseError(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode < 300 {
		return nil
	}
	return errors.New(r.Status)
}
