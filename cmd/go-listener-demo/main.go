package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lucasew/go-getlistener"
)

func main() {
	ln, err := getlistener.GetListener()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("serving on %s", ln.Addr())

	err = http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "works!")
	}))
	if err != nil {
		log.Fatal(err)
	}
}
