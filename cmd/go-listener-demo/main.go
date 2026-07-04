package main

import (
	"fmt"
	"net/http"

	"github.com/lucasew/go-getlistener"
	"github.com/lucasew/go-getlistener/internal/errorhandler"
)

type Server struct{}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "works!")
}

func main() {
	ln, err := getlistener.GetListener()
	if err != nil {
		errorhandler.ReportErrorAndExit(err, "Failed to set up listener")
	}
	err = http.Serve(ln, Server{})
	if err != nil {
		errorhandler.ReportErrorAndExit(err, "Failed to set up HTTP server")
	}

}
