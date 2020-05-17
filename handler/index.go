package handler

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/sah4ez/pspk/pkg/pspk"
	"github.com/sah4ez/pspk/pkg/validation"
	qrcode "github.com/skip2/go-qrcode"
)

type pub []byte

const (
	maxLimit = 500
)

type Request struct {
	ID     bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string        `json:"name,omitempty" bson:"name"`
	Key    string        `json:"key,omitempty" bson:"key"`
	Method string        `json:"method,omitempty" bson:"-"`
	Data   string        `json:"data,omitempty" bson:"-"`
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
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Origin, X-Requested-With, Content-Type, Accept")

	if r.Method == "OPTIONS" {
		return
	}
	once.Do(func() { initConnection(w, resp) })

	value := r.Header.Get("X-Access-Token")

	if r.Method == http.MethodGet {
		query := r.URL.Query()
		if query.Get(pspk.LinkKey) != "" {
			if err := GetByLink(w, r); err != nil {
				resp["error"] = err.Error()
				//todo: fixed handled error
				_ = json.NewEncoder(w).Encode(resp)
			}
			return
		}
		switch query.Get(pspk.OutputKey) {
		case "json", "json-map":
			w.Header().Set("Content-Type", "application/json")
			if err := GetKeysInJson(w, r); err != nil {
				resp["error"] = err.Error()
				//todo: fixed handled error
				_ = json.NewEncoder(w).Encode(resp)
			}
			return
		case "json-array":
			w.Header().Set("Content-Type", "application/json")
			if err := GetKeysInJsonArray(w, r); err != nil {
				resp["error"] = err.Error()
				//todo: fixed handled error
				_ = json.NewEncoder(w).Encode(resp)
			}
			return
		}

		if name := query.Get(pspk.QRCodeKey); name != "" {
			var key Request
			key, err = ByName(name)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				resp["error"] = err.Error()
				//todo: fixed handled error
				_ = json.NewEncoder(w).Encode(resp)
				return
			}

			var png []byte
			png, err := qrcode.Encode(key.Key, qrcode.Highest, 256)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				resp["error"] = err.Error()
				//todo: fixed handled error
				_ = json.NewEncoder(w).Encode(resp)
				return
			}

			w.Header().Set("Content-Type", "image/png")
			if _, err := io.Copy(w, bytes.NewReader(png)); err != nil {
				resp["error"] = err.Error()
				//todo: fixed handled error
				_ = json.NewEncoder(w).Encode(resp)
				return
			}
			return
		}
		if name := query.Get(pspk.NameKey); name != "" {
			w.Header().Set("Content-Type", "application/json")
			var key Request
			key, err = ByName(name)
			if err != nil {
				resp["error"] = err.Error()
				//todo: fixed handled error
				_ = json.NewEncoder(w).Encode(resp)
				return
			}
			if err = json.NewEncoder(w).Encode([]Request{key}); err != nil {
				resp["error"] = err.Error()
				//todo: fixed handled error
				_ = json.NewEncoder(w).Encode(resp)
				return
			}
			return
		}
		if name := query.Get(pspk.NameSearchKey); name != "" {
			w.Header().Set("Content-Type", "application/json")
			keys, err := FindByName(name)
			if err != nil {
				resp["error"] = err.Error()
				//todo: fixed handled error
				_ = json.NewEncoder(w).Encode(resp)
				return
			}
			if err = json.NewEncoder(w).Encode(keys); err != nil {
				resp["error"] = err.Error()
				//todo: fixed handled error
				_ = json.NewEncoder(w).Encode(resp)
				return
			}
			return
		}
		if err := Get(w, r); err != nil {
			resp["error"] = err.Error()
			//todo: fixed handled error
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
	}
	resp["access"] = value == token

	defer func() {
		if err != nil {
			resp["error"] = err.Error()
		}
		//todo: fixed handled error
		_ = json.NewEncoder(w).Encode(resp)
	}()

	if r.Method == http.MethodPost {
		var keyRequest Request

		err = json.NewDecoder(r.Body).Decode(&keyRequest)
		if err != nil {
			return
		}

		switch strings.ToLower(keyRequest.Method) {
		case pspk.LinkKey:
			if keyRequest.Data == "" {
				err = fmt.Errorf("empty data")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if len(keyRequest.Data) > 2048 {
				err = fmt.Errorf("very large data, max 2048 symbols")
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

func LoadArray(lastID bson.ObjectId, limit int) (keys []Request, err error) {
	sess := session.Copy()
	defer sess.Close()

	c := sess.DB("pspk").C("keys")

	query := bson.M{}
	if lastID != "" {
		query["_id"] = bson.M{"$gt": lastID}
	}
	err = c.Find(query).Limit(limit).All(&keys)
	if err != nil {
		return
	}

	return
}

func Load(lastID bson.ObjectId, limit int) (keys map[string]pub, err error) {
	keys = map[string]pub{}
	sess := session.Copy()
	defer sess.Close()

	c := sess.DB("pspk").C("keys")

	result := make([]Request, 0, 64)

	query := bson.M{}
	if lastID != "" {
		query["_id"] = bson.M{"$gt": lastID}
	}
	err = c.Find(query).Limit(limit).All(&result)
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

func FindByName(name string) (result []Request, err error) {
	sess := session.Copy()
	defer sess.Close()

	c := sess.DB("pspk").C("keys")

	err = c.Find(bson.M{"name": bson.RegEx{Pattern: name + ".*"}}).Sort("name").All(&result)
	if err != nil {
		return
	}

	return
}

func initConnection(w io.Writer, resp map[string]interface{}) {
	var err error
	dialInfo, err := mgo.ParseURL(addr)
	if err != nil {
		resp["error"] = err.Error()
		resp["cause"] = "parse"
		resp["url"] = addr
		//todo: fixed handled error
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	dialInfo.Timeout = 5 * time.Second

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		tlsConfig := &tls.Config{}

		if local != nil && *local {
			cer, err := tls.LoadX509KeyPair("test_certs/localhost.crt", "test_certs/localhost.key")
			if err != nil {
				return nil, fmt.Errorf("load key pair failed: %w", err)
			}
			// Load CA cert
			caCert, err := ioutil.ReadFile("test_certs/rootCA.crt")
			if err != nil {
				return nil, fmt.Errorf("read root ca failed: %w", err)
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)

			tlsConfig.Certificates = []tls.Certificate{cer}
			tlsConfig.RootCAs = caCertPool
			tlsConfig.ServerName = "localhost"
			tlsConfig.BuildNameToCertificate()
		}

		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)

		if err != nil {
			_, _ = fmt.Fprint(output, err.Error())
			resp["error"] = err.Error()
			resp["cause"] = "dial func"
			//todo: fixed handled error
			_ = json.NewEncoder(w).Encode(resp)
		}
		return conn, err
	}

	session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		resp["error"] = err.Error()
		resp["cause"] = "dial"
		//todo: fixed handled error
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	c := session.DB("pspk").C("keys")
	err = c.EnsureIndex(mgo.Index{Key: []string{"name"}, Unique: true})
	if err != nil {
		resp["error"] = err.Error()
		resp["cause"] = "create index"
		//todo: fixed handled error
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	c = session.DB("pspk").C("links")
	err = c.EnsureIndex(mgo.Index{Key: []string{"create_at"}, ExpireAfter: 60 * 60 * 24 * time.Second})
	if err != nil {
		resp["error"] = err.Error()
		resp["cause"] = "create index"
		//todo: fixed handled error
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
}
func GetByLink(w io.Writer, r *http.Request) (err error) {
	query := r.URL.Query()
	id := query.Get(pspk.LinkKey)
	data, err := FindByLinkId(id)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(map[string]string{"data": data})
	return
}

func loadPagination(r *http.Request) (id bson.ObjectId, limit int) {
	query := r.URL.Query()

	if keyID := query.Get(pspk.LastIDKEy); keyID != "" {
		id = bson.ObjectIdHex(keyID)
	}
	limitStr := query.Get(pspk.LimitKey)

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = maxLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	return
}

func GetKeysInJson(w io.Writer, r *http.Request) (err error) {
	lastID, limit := loadPagination(r)
	data, err := Load(lastID, limit)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(data)
	return
}

func GetKeysInJsonArray(w io.Writer, r *http.Request) (err error) {
	lastID, limit := loadPagination(r)
	data, err := LoadArray(lastID, limit)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(data)
	return
}

func Get(w io.Writer, _ *http.Request) (err error) {
	tmpl, err := template.New("index").Funcs(template.FuncMap{
		"base64": func(b []byte) string {
			return base64.StdEncoding.EncodeToString(b)
		},
	}).Parse(body)
	if err != nil {
		return fmt.Errorf("parse body: %s", err.Error())
	}

	var keys map[string]pub

	keys, err = Load("", maxLimit)
	if err != nil {
		return fmt.Errorf("load all key: %s", err.Error())
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
