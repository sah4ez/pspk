package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sah4ez/pspk/handler"
)

var (
	//Version current tools
	Version string
	// Hash revision number from git
	Hash string
	// BuildDate when building this utilities
	BuildDate string

	output   io.Writer = os.Stdout
	addrFlag           = flag.String("addr", "0.0.0.0:8080", "bind addr for HTTP")
)

func init() {

	flag.Parse()
	if *addrFlag == "" {
		panic("invalid addr")
	}
}

func main() {
	fmt.Fprintln(output, "init server")
	fmt.Fprintln(output, "version", Version)
	fmt.Fprintln(output, "Hash", Hash)
	fmt.Fprintln(output, "BuildDate", BuildDate)

	http.HandleFunc("/", handler.Handler)

	errs := make(chan error)

	if err := http.ListenAndServe(*addrFlag, nil); err != nil {
		os.Exit(1)
	}
	err := <-errs
	if err != nil {
		fmt.Fprintln(output, err.Error())
		os.Exit(1)
	}
}
