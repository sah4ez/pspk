// +build js,wasm

package main

import (
	"fmt"

	"github.com/sah4ez/pspk/pkg/wasm"
)

func main() {
	c := make(chan struct{}, 0)

	fmt.Println("wasm module loaded")

	wasm.Load()

	<-c
}
