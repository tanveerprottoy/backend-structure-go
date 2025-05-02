package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Option func(*Server)

// WithReadTimeout sets the ReadTimeout for the server
func WithReadTimeout(t time.Duration) Option {
	return func(srv *Server) {
		srv.httpServer.ReadTimeout = t
	}
}

// WithReadHeaderTimeout sets the ReadHeaderTimeout for the server
func WithReadHeaderTimeout(t time.Duration) Option {
	return func(srv *Server) {
		srv.httpServer.ReadHeaderTimeout = t
	}
}

// WithWriteTimeout sets the WriteTimeout for the server
func WithWriteTimeout(t time.Duration) Option {
	return func(srv *Server) {
		srv.httpServer.WriteTimeout = t
	}
}

type Server struct {
	httpServer *http.Server
	// empty struct consumes zero memory
	// this channel is used to wait for idle connections to be closed
	// before shutting down the server
	idleConnsClosed chan struct{}
}

// NewServer initializes the server
func NewServer(address string, handler http.Handler, opts ...Option) *Server {
	srv := &Server{
		httpServer: &http.Server{
			Addr:    address,
			Handler: handler,
		},
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (s *Server) Shutdown() {
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server shutdown error: %v", err)
	}
}

// ConfigureGracefulShutdown configures graceful shutdown
func (s *Server) ConfigureGracefulShutdown(defferedFunc func()) {
	// code to support graceful shutdown
	s.idleConnsClosed = make(chan struct{})

	go func() {
		// this func listens for SIGINT and initiates
		// shutdown when SIGINT is received
		ch := make(chan os.Signal, 1)

		// register ch to receive interrupt signal
		signal.Notify(ch, os.Interrupt)

		// this will block, wait for signal
		// receive data from ch
		<-ch

		// Received an interrupt signal, shut down.
		log.Printf("Received an interrupt signal")

		if defferedFunc != nil {
			defer defferedFunc()
		}

		s.Shutdown()

		// close the idle connection close channel
		close(s.idleConnsClosed)
	}()
}

// Start starts the server
func (s *Server) Start() {
	log.Println("server starting")

	// if err == http.ErrServerClosed do nothing
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	// wait for idle connections to be closed
	<-s.idleConnsClosed

	log.Println("server shutdown")
}

// HTTPServer returns the http server
func (s *Server) HTTPServer() *http.Server {
	return s.httpServer
}
