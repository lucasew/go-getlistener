// Package main provides a demo server for the go-getlistener library.
package main

import (
	"fmt"
	"net/http"

	"github.com/lucasew/go-getlistener"
)

// Server implements the http.Handler interface.
type Server struct{}

func (s Server) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "works!")
}

func main() {
	ln, err := getlistener.GetListener()
	if err != nil {
		panic(err)
	}
	err = http.Serve(ln, Server{})
	if err != nil {
		panic(err)
	}

}
