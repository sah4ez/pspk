// +build js,wasm

package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"syscall/js"

	"github.com/sah4ez/pspk/pkg/pspk"
	"github.com/sah4ez/pspk/pkg/utils"
	"github.com/skip2/go-qrcode"
)

func main() {
	fmt.Println("wasm module loaded")

	api := pspk.New("http://127.0.0.1:8080/")
	fs := utils.NewWasmStorage()
	cli := pspk.NewPSPKcli(api, nil, "/", "https://127.0.0.1:8080", os.Stdout, fs)

	var (
		name   string
		qrCode bool
	)

	name = js.Global().Get("pub_name").Get("value").String()
	fmt.Println(name)
	qrCode = js.Global().Get("qr_codes").Get("checked").Bool()

	var result error
	result = cli.Publish(name)

	if result != nil {
		js.Global().Set("PublishError", result.Error())
		return
	}

	kPub, _ := fs.Read("/"+name, "pub.bin")
	js.Global().Get("pub_key").Set("value", base64.StdEncoding.EncodeToString(kPub))

	kPriv, _ := fs.Read("/"+name, "key.bin")
	js.Global().Get("priv_key").Set("value", base64.StdEncoding.EncodeToString(kPriv))

	if qrCode {
		pubCanvas, pubDone := loadCanvas("qr_pub")
		drawImage(kPub, pubCanvas, pubDone)

		privCanvas, privDone := loadCanvas("qr_priv")
		drawImage(kPriv, privCanvas, privDone)
		<-pubDone
		<-privDone
	}

	js.Global().Set("PublishError", "")
}

func loadCanvas(name string) (js.Value, chan struct{}) {
	canvas := js.Global().Get("document").Call("getElementById", name)
	canvas.Set("width", js.ValueOf(256))
	canvas.Set("height", js.ValueOf(256))
	done := make(chan struct{})
	return canvas, done
}

func drawImage(data []byte, canvas js.Value, done chan struct{}) {
	k, err := qrcode.Encode(string(data[:]), qrcode.Medium, 256)
	if err != nil {
		fmt.Println("err", err.Error())
	}
	image := js.Global().Call("eval", "new Image()")
	image.Set("src", "data:image/png;base64,"+base64.StdEncoding.EncodeToString(k))
	js.Global().Set("qr_code_img", image)

	ctx := canvas.Call("getContext", "2d")
	ctx.Call("clearRect", 0, 0, 256, 256)
	js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ctx.Call("drawImage", image, 0, 0)
		close(done)
		return nil
	}), js.ValueOf(1000))
}
