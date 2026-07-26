package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lucasew/go-getlistener"
)

func main() {
	ln, err := getlistener.GetListener()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("serving on %s", ln.Addr())

	// systemd and interactive runs both send SIGTERM/SIGINT for a clean stop.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "works!")
		}),
		// Bound request/response lifetimes so slow or idle clients cannot hold
		// connections open forever (ReadHeaderTimeout alone is not enough).
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		// Cancel r.Context() on SIGTERM/SIGINT so handlers that respect
		// request context can exit before Shutdown's deadline.
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	errCh := make(chan error, 1)
	go func() {
		err := srv.Serve(ln)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
		if err := <-errCh; err != nil {
			log.Fatal(err)
		}
	case err := <-errCh:
		if err != nil {
			log.Fatal(err)
		}
	}
}
