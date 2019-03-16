package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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

func Handler(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		resp = make(map[string]interface{})
	)

	defer func() {
		if err != nil {
			resp["error"] = err.Error()
		}
		json.NewEncoder(w).Encode(resp)
	}()

	value := r.Header.Get("X-Access-Token")

	resp["access"] = value == token

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
