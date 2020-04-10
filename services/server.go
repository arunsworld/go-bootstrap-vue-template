package services

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ServerConfig configures a server
type ServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func (config *ServerConfig) updateWithDefaults() {
	if config.ReadTimeout == time.Duration(0) {
		config.ReadTimeout = 5 * time.Second
	}
	if config.WriteTimeout == time.Duration(0) {
		config.WriteTimeout = 10 * time.Second
	}
	if config.IdleTimeout == time.Duration(0) {
		config.IdleTimeout = 120 * time.Second
	}
}

// ServerStart starts up the server
func ServerStart(config ServerConfig, srvMux http.Handler) {
	config.updateWithDefaults()
	srv := &http.Server{
		Addr:         config.Addr,
		Handler:      srvMux,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}
	ctx := trapKillSignalForGracefulShutdown()
	go func() {
		<-ctx.Done()
		log.Println("shutting down HTTP server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Println(err)
		}
		log.Println("HTTP server shutdown complete...")
	}()
	log.Printf("starting HTTP server on: %s", config.Addr)
	srv.ListenAndServe()
	// srv.ListenAndServeTLS("server.crt", "server.key")
}

func trapKillSignalForGracefulShutdown() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGTERM)

		<-signals

		log.Println("KILL signal received, initiating termination...")

		cancel()
	}()

	return ctx
}
