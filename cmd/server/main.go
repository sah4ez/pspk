package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
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

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT)

	errs := make(chan error)

	go func() {
		fmt.Fprintln(output, "started server on", *addrFlag)
		err := http.ListenAndServe(*addrFlag, nil)
		errs <- errors.Wrap(err, "Failed listening server")
	}()

	go func() {
		if err := <-errs; err != nil {
			fmt.Fprintln(output, err.Error())
			os.Exit(1)
		}
	}()

	<-shutdown
	fmt.Fprintln(output, "Server terminated")
}
