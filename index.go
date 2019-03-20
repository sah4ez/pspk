package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type pub []byte

type Request struct {
	Name string `json:"name" bson:"name"`
	Key  string `json:"key,omitempty" bson:"key"`
}

var (
	token   = os.Getenv("ACCESS_TOKEN")
	user    = os.Getenv("DB_USER")
	pass    = os.Getenv("DB_PASS")
	hosts   = os.Getenv("DB_HOSTS")
	addr    = fmt.Sprintf("mongodb://%s:%s@%s/admin", user, pass, hosts)
	session *mgo.Session
	once    = &sync.Once{}

	errKeyNotFound = fmt.Errorf("key not found")
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

	once.Do(func() { initConnection(w, resp) })

	value := r.Header.Get("X-Access-Token")

	resp["access"] = value == token

	if r.Method == http.MethodGet {
		tmpl, err := template.New("index").Funcs(template.FuncMap{
			"base64": func(b []byte) string {
				return base64.StdEncoding.EncodeToString(b)
			},
		}).Parse(body)
		if err != nil {
			resp["error"] = err.Error()
			return
		}

		keys, err := Load()
		if err != nil {
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
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

		if keyRequest.Key == "" {
			key, err := ByName(keyRequest.Name)
			if err != nil {
				if err == errKeyNotFound {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			resp["key"] = key.Key
			return
		}

		err = PutKey(keyRequest)
		if mgo.IsDup(err) {
			err = fmt.Errorf("name %s is exist", keyRequest.Name)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp["msg"] = "added"

		w.WriteHeader(http.StatusCreated)
		return
	}
}

func PutKey(p Request) (err error) {
	sess := session.Copy()
	defer sess.Close()

	c := sess.DB("pspk").C("keys")

	var b pub
	b, err = base64.StdEncoding.DecodeString(p.Key)
	if err != nil {
		return
	}

	if len(b) != 32 {
		return fmt.Errorf("should be 32-bytes key")
	}

	return c.Insert(&p)
}

func ByName(name string) (p Request, err error) {
	sess := session.Copy()
	defer sess.Close()

	c := sess.DB("pspk").C("keys")

	p = Request{}

	err = c.Find(bson.M{"name": name}).One(&p)
	if err != nil {
		if err == mgo.ErrNotFound {
			return Request{}, errKeyNotFound
		}
		return
	}

	return
}

func Load() (keys map[string]pub, err error) {
	keys = map[string]pub{}
	sess := session.Copy()
	defer sess.Close()

	c := sess.DB("pspk").C("keys")

	result := []Request{}

	err = c.Find(bson.M{}).All(&result)
	if err != nil {
		return
	}

	for _, k := range result {
		b, err := base64.StdEncoding.DecodeString(k.Key)
		if err != nil {
			keys[k.Name] = pub{}
			continue
		}
		keys[k.Name] = b
	}

	return
}

func initConnection(w io.Writer, resp map[string]interface{}) {
	dialInfo, err := mgo.ParseURL(addr)
	if err != nil {
		resp["error"] = err.Error()
		resp["cause"] = "parse"
		json.NewEncoder(w).Encode(resp)
		return
	}

	dialInfo.Timeout = 5 * time.Second

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		tlsConfig := &tls.Config{}
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)

		if err != nil {
			resp["error"] = err.Error()
			resp["cause"] = "dial func"
			json.NewEncoder(w).Encode(resp)
		}
		return conn, err
	}

	session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		resp["error"] = err.Error()
		resp["cause"] = "dial"
		json.NewEncoder(w).Encode(resp)
		return
	}

	c := session.DB("pspk").C("keys")
	err = c.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true})
	if err != nil {
		resp["error"] = err.Error()
		resp["cause"] = "create index"
		json.NewEncoder(w).Encode(resp)
		return
	}

}
