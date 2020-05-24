// +build js,wasm

package main

import (
	"encoding/base64"
	"fmt"
	"syscall/js"

	"github.com/sah4ez/pspk/pkg/pspk"
	"github.com/sah4ez/pspk/pkg/utils"
)

var (
	ednpoint = "http://127.0.0.1:8080"
)

func main() {
	fmt.Println("wasm module loaded")

	api := pspk.New(ednpoint)
	fs := utils.NewWasmStorage()

	out := utils.NewMessageWriter()
	cli := pspk.NewPSPKcli(api, nil, "/", ednpoint, out, fs)

	var (
		keyBase64  string
		dataBase64 string
	)

	keyBase64 = js.Global().Get("private_key").Get("value").String()
	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		js.Global().Set("PublishError", err.Error())
		return
	}
	dataBase64 = js.Global().Get("text_dec").Get("value").String()

	fs.Write("/web", "key.bin", key)

	err = cli.EphemeralDecrypt("web", dataBase64)

	if err != nil {
		js.Global().Set("PublishError", err.Error())
		fmt.Println(err.Error())
		return
	}
	fmt.Println(out.Read())

	js.Global().Get("copy_dec").Set("value", out.Read())
	js.Global().Get("copy_dec").Call("select")
	js.Global().Get("document").Call("execCommand", "copy")
	js.Global().Set("PublishError", "")
}
