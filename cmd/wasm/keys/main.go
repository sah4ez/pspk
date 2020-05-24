// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/sah4ez/pspk/pkg/pspk"
)

func main() {
	fmt.Println("wasm module loaded")

	var name string

	name = js.Global().Get("pub_name").Get("value").String()

	api := pspk.New("http://127.0.0.1:8080/?name_regex=" + name)
	options := pspk.GetAllOptions{}
	result := api.GetAll(options)

	if result.Error != nil {
		js.Global().Set("PublishError", result.Error.Error())
		fmt.Println("result", result.Error.Error())
		return
	}

	bytes, err := json.Marshal(result.Keys)
	if err != nil {
		js.Global().Set("PublishError", err.Error())
		fmt.Println("marshal", err.Error())
		return
	}
	fmt.Println(string(bytes))
	// js.Global().Get("table").Set("value", string(bytes))

	js.Global().Set("PublishError", "")
}
