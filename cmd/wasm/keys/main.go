// +build js,wasm

package main

import (
	"fmt"
	"strings"
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

	tables := js.Global().Get("document").Call("getElementsByClassName", "table")
	tbody := tables.Get("0").Get("lastElementChild")

	var row = `
<tr>
	<th scope="row">%s</th>
	<td>%s</td>
	<td>%s</td>
</tr>
	`
	rows := make([]string, len(result.Keys))

	for i, key := range result.Keys {
		rows[i] = fmt.Sprintf(row, key.ID, key.Name, key.Key)
	}
	tbody.Set("innerHTML", strings.Join(rows, "\n"))

	js.Global().Set("PublishError", "")
}
