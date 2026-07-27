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
	// log.Fatal calls os.Exit and skips defers; keep it only after run returns
	// so signal.NotifyContext stop and shutdown cancel always run.
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ln, err := getlistener.GetListener()
	if err != nil {
		return err
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
			return err
		}
		return <-errCh
	case err := <-errCh:
		return err
	}
}
