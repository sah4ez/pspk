package handler

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/sah4ez/pspk/pkg/validation"
)

type pub []byte

const (
	NameKey = "name_key"
	LinkKey = "link"
)

type Request struct {
	Name   string `json:"name,omitempty" bson:"name"`
	Key    string `json:"key,omitempty" bson:"key"`
	Method string `json:"method,omitempty" bson:"-"`
	Data   string `json:"data,omitempty" bson:"-"`
}

var (
	token   = os.Getenv("ACCESS_TOKEN")
	user    = os.Getenv("DB_USER")
	pass    = os.Getenv("DB_PASS")
	hosts   = os.Getenv("DB_HOSTS")
	local   = flag.Bool("local", false, "run with connect to local mongo")
	addr    = fmt.Sprintf("mongodb://%s:%s@%s/admin", user, pass, hosts)
	session *mgo.Session
	once    = &sync.Once{}

	output io.Writer = os.Stdout

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

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		return
	}
	once.Do(func() { initConnection(w, resp) })

	value := r.Header.Get("X-Access-Token")

	resp["access"] = value == token

	if r.Method == http.MethodGet {
		idLink := r.URL.Query().Get(LinkKey)
		if idLink != "" {
			if err := GetByLink(w, r); err != nil {
				resp["error"] = err.Error()
				json.NewEncoder(w).Encode(resp)
			}
			return
		}
		if err := Get(w, r); err != nil {
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}
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

		switch strings.ToLower(keyRequest.Method) {
		case LinkKey:
			if keyRequest.Data == "" {
				err = fmt.Errorf("empty data")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if len(keyRequest.Data) > 2048 {
				err = fmt.Errorf("Very large data, max 2048 symbols")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			id, err := GenerateLinkId(keyRequest.Data)
			if err != nil {
				return
			}
			resp["link"] = id
			return
		default:
			if keyRequest.Key == "" && keyRequest.Name == "" {
				err = fmt.Errorf("not set values")
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err = validation.CheckLimitNameLen(keyRequest.Name); err != nil {
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
		}

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

func decode(result []Request) map[string]pub {
	keys := map[string]pub{}
	for _, k := range result {
		b, err := base64.StdEncoding.DecodeString(k.Key)
		if err != nil {
			keys[k.Name] = pub{}
			continue
		}
		keys[k.Name] = b
	}
	return keys
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

	keys = decode(result)
	return
}

func FindByLinkId(id string) (data string, err error) {
	sess := session.Copy()
	defer sess.Close()

	c := sess.DB("pspk").C("links")

	var body struct {
		Data string `bson:"data"`
	}

	err = c.FindId(bson.ObjectIdHex(id)).One(&body)
	return body.Data, err
}

func GenerateLinkId(data string) (id string, err error) {
	sess := session.Copy()
	defer sess.Close()

	c := sess.DB("pspk").C("links")

	var link = struct {
		ID        bson.ObjectId `bson:"_id"`
		Data      string        `bson:"data"`
		CreatedAt time.Time     `bson:"created_at"`
	}{
		ID:        bson.NewObjectId(),
		Data:      data,
		CreatedAt: time.Now(),
	}

	err = c.Insert(&link)
	return link.ID.Hex(), err
}

func FindByName(name string) (keys map[string]pub, err error) {
	sess := session.Copy()
	defer sess.Close()

	c := sess.DB("pspk").C("keys")

	result := []Request{}

	err = c.Find(bson.M{"name": "/.*" + name + ".*/"}).All(&result)
	if err != nil {
		return
	}

	keys = decode(result)
	return
}

func GetByLink(w io.Writer, r *http.Request) (err error) {
	query := r.URL.Query()
	id := query.Get(LinkKey)
	data, err := FindByLinkId(id)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(map[string]string{"data": data})
	return
}

func Get(w io.Writer, r *http.Request) (err error) {
	tmpl, err := template.New("index").Funcs(template.FuncMap{
		"base64": func(b []byte) string {
			return base64.StdEncoding.EncodeToString(b)
		},
	}).Parse(body)
	if err != nil {
		return fmt.Errorf("parse body: %s", err.Error())
	}

	var keys map[string]pub

	query := r.URL.Query()
	name := query.Get(NameKey)
	if name != "" {
		keys, err = FindByName(name)
		if err != nil {
			return fmt.Errorf("load key by name: %s", err.Error())
		}

	} else {
		keys, err = Load()
		if err != nil {
			return fmt.Errorf("load all key: %s", err.Error())
		}
	}

	err = tmpl.Execute(w, PspkStore{
		Title: "PSPK kv store",
		Keys:  keys,
	})
	if err != nil {
		return fmt.Errorf("execute template: %s", err.Error())
	}
	return
}
