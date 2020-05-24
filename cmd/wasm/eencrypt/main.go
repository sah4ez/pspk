// +build js,wasm

package main

import (
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
		pubName     string
		message     string
		withoutLink = false
	)

	pubName = js.Global().Get("pub_name").Get("value").String()
	message = js.Global().Get("text_enc").Get("value").String()

	err := cli.EphemeralEncrypt(message, pubName, withoutLink)

	if err != nil {
		js.Global().Set("PublishError", err.Error())
		fmt.Println(err.Error())
		return
	}

	fmt.Println(out.Read())

	js.Global().Get("copy_enc").Set("value", out.Read())
	js.Global().Get("copy_enc").Call("select")
	js.Global().Get("document").Call("execCommand", "copy")
	js.Global().Set("PublishError", "")
}
