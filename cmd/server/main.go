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
		fmt.Println(output, "invalid flag '-addr', please enter addr like ':8080'")
		os.Exit(1)
	}
}

func main() {
	fmt.Fprintln(output, "init server")
	fmt.Fprintln(output, "version", Version)
	fmt.Fprintln(output, "Hash", Hash)
	fmt.Fprintln(output, "BuildDate", BuildDate)

	http.HandleFunc("/", handler.Handler)

	errs := make(chan error)

	errs <- http.ListenAndServe(*addrFlag, nil)

	err := <-errs
	if err != nil {
		fmt.Fprintln(output, err.Error())
		os.Exit(1)
	}
}
