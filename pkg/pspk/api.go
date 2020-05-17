package pspk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

type Link struct {
	Link  string
	Error error
}

type PublicKey struct {
	Key   []byte
	Error error
}

type DownloadData struct {
	Data  string
	Error error
}

type Keys struct {
	Keys  []Key
	Error error
}

type PSPK interface {
	Publish(name string, key []byte) (err error)
	Load(name string) (key *PublicKey)
	GenerateLink(data string) (link *Link)
	DownloadByLink(link string) (data *DownloadData)
	GetAll(opts GetAllOptions) (keys *Keys)
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

func (p *pspk) Publish(name string, key []byte) (err error) {
	if err = validation.CheckLimitNameLen(name); err != nil {
		return
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
		return
	}

	resource := struct {
		Access bool   `json:"access"`
		Msg    string `json:"msg"`
	}{}

	err = p.doRequest(req, &resource)
	if err != nil {
		err = fmt.Errorf("Do requset laod failed: %s", err)
		return
	}
	return
}

func (p *pspk) Load(name string) (key *PublicKey) {

	key = &PublicKey{}
	if err := validation.CheckLimitNameLen(name); err != nil {
		key.Error = fmt.Errorf("lager limt of name: %s", err)
		return
	}
	body := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}
	req, err := p.newRequest("POST", body, nil)
	if err != nil {
		key.Error = fmt.Errorf("Create requset load failed: %s", err)
		return
	}
	resource := struct {
		Key string `json:"key"`
	}{}
	err = p.doRequest(req, &resource)
	if err != nil {
		key.Error = fmt.Errorf("Do requset laod failed: %s", err)
		return
	}
	key.Key, key.Error = base64.StdEncoding.DecodeString(resource.Key)
	return
}

func (p *pspk) GenerateLink(data string) (link *Link) {

	link = &Link{}
	if len(data) == 0 {
		link.Error = fmt.Errorf("empty data for generation link")
		return
	}
	body := struct {
		Method string `json:"method"`
		Data   string `json:"data"`
	}{
		Method: LinkKey,
		Data:   data,
	}
	req, err := p.newRequest("POST", body, nil)
	if err != nil {
		link.Error = fmt.Errorf("Create request for generate link: %s", err)
		return
	}
	resource := struct {
		Link string `json:"link"`
	}{}
	err = p.doRequest(req, &resource)
	if err != nil {
		link.Error = fmt.Errorf("Generate link through server failed: %s", err)
		return
	}
	link.Link = resource.Link
	return
}

func (p *pspk) DownloadByLink(link string) (data *DownloadData) {
	data = &DownloadData{}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		data.Error = err
		return
	}
	resp, err := p.client.Do(req)
	if err != nil {
		data.Error = err
		return
	}
	defer resp.Body.Close()

	b := &struct {
		Data string `json:"data"`
	}{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&b)
	if err != nil {
		data.Error = err
		return
	}
	data.Data = b.Data
	return
}

func (p *pspk) GetAll(options GetAllOptions) (keys *Keys) {

	keys = &Keys{}
	req, err := p.newRequest("GET", nil, options)
	if err != nil {
		keys.Error = err
		return
	}
	var resource []Key
	err = p.doRequest(req, &resource)
	if err != nil {
		keys.Error = err
		return
	}
	keys.Keys = resource
	return
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
