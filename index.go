package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sync"
)

type pub []byte

type Request struct {
	Name string `json:"name"`
	Key  string `json:"key,omitempty"`
}

var (
	token = os.Getenv("ACCESS_TOKEN")

	lock = sync.RWMutex{}
	keys = make(map[string]pub)
)

type PspkStore struct {
	Title string
	Keys  map[string]pub
}

var body = `<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Title}}</title>
</head>
<body>
	{{range $k, $v := .Keys}}
	<div>key: {{$k}}</div>
	<div>val: {{base64 $v}}</div>
	{{end}}
</body>
</html>`

func Handler(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		resp = make(map[string]interface{})
	)

	value := r.Header.Get("X-Access-Token")

	resp["access"] = value == token

	if r.Method == http.MethodGet {
		lock.RLock()
		defer lock.RUnlock()
		tmpl, err := template.New("index").Funcs(template.FuncMap{
			"base64": func(b []byte) string {
				return base64.StdEncoding.EncodeToString(b)
			},
		}).Parse(body)
		if err != nil {
			resp["error"] = err.Error()
		}
		err = tmpl.Execute(w, PspkStore{
			Title: "PSPK kv store",
			Keys:  keys,
		})
		if err != nil {
			resp["error"] = err.Error()
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	defer func() {
		if err != nil {
			resp["error"] = err.Error()
		}
		json.NewEncoder(w).Encode(resp)
	}()

	if r.Method == http.MethodPost {
		var keyRequest Request

		err = json.NewDecoder(r.Body).Decode(&keyRequest)
		if err != nil {
			return
		}

		if keyRequest.Key == "" && keyRequest.Name == "" {
			err = fmt.Errorf("not set values")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		lock.Lock()
		defer lock.Unlock()

		if keyRequest.Key == "" {
			key, ok := keys[keyRequest.Name]
			if ok {
				resp["key"] = base64.StdEncoding.EncodeToString(key)
			} else {
				err = fmt.Errorf("name %s not found", keyRequest.Name)
				w.WriteHeader(http.StatusNotFound)
				return
			}
		} else {
			_, ok := keys[keyRequest.Name]
			if ok {
				err = fmt.Errorf("name %s is exist", keyRequest.Name)
				w.WriteHeader(http.StatusBadRequest)
				return
			} else {
				var b pub
				b, err = base64.StdEncoding.DecodeString(keyRequest.Key)
				if err != nil {
					return
				}

				if len(b) != 32 {
					err = fmt.Errorf("should be 32-bytes key")
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				keys[keyRequest.Name] = b
				resp["msg"] = "added"

				w.WriteHeader(http.StatusCreated)
				return
			}
		}
	}
}
