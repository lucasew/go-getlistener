// Package main implements a demo server.
package main

import (
	"fmt"
	"net/http"

	"github.com/lucasew/go-getlistener"
)

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
