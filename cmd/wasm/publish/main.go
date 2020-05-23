// +build js,wasm

package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"syscall/js"

	"github.com/sah4ez/pspk/pkg/pspk"
	"github.com/sah4ez/pspk/pkg/utils"
)

func main() {
	fmt.Println("wasm module loaded")

	api := pspk.New("http://127.0.0.1:8080/")
	fs := utils.NewWasmStorage()
	cli := pspk.NewPSPKcli(api, nil, "/", "http://127.0.0.1:8080", os.Stdout, fs)

	var name string

	name = js.Global().Get("publish_name").String()
	fmt.Println(name)

	var result error
	result = cli.Publish(name)

	if result != nil {
		js.Global().Set("PublishError", result.Error())
		return
	}

	k, _ := fs.Read("/"+name, "pub.bin")
	js.Global().Set("pub_key", base64.StdEncoding.EncodeToString(k))

	k, _ = fs.Read("/"+name, "key.bin")
	js.Global().Set("priv_key", base64.StdEncoding.EncodeToString(k))

	js.Global().Set("PublishError", "")
}
