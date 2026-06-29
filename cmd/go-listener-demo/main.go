package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/lucasew/go-getlistener"
)

type Server struct{}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "works!")
}

func main() {
	ln, err := getlistener.GetListener()
	if err != nil {
		getlistener.ReportError(err, "failed to get listener")
		os.Exit(1)
	}
	err = http.Serve(ln, Server{})
	if err != nil {
		getlistener.ReportError(err, "server failed")
		os.Exit(1)
	}

}
