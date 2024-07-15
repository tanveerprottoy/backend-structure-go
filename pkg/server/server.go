package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type Server struct {
	httpServer *http.Server
	// empty struct consumes zero memory
	idleConnsClosed chan struct{}
}

// NewServer initializes the server
func NewServer(address string, handler http.Handler, opts ...any) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    address,
			Handler: handler,
			/* ReadTimeout:       constant.ServerReadTimeout * time.Second,
			ReadHeaderTimeout: constant.ServerReadHeaderTimeout * time.Second,
			WriteTimeout:      constant.ServerWriteTimeout * time.Second,
			*/
		},
	}
}

func (s *Server) Shutdown() {
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP Server shutdown error: %v", err)
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
		signal.Notify(ch, os.Interrupt)
		<-ch
		// Received an interrupt signal, shut down.
		log.Printf("Received an interrupt signal")
		if defferedFunc != nil {
			defer defferedFunc()
		}
		s.Shutdown()
		close(s.idleConnsClosed)
	}()
}

// Start starts the server
func (s *Server) Start() {
	log.Println("Server starting")
	// if err == http.ErrServerClosed do nothing
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP Server ListenAndServe: %v", err)
	}
	// wait for idle connections to be closed
	<-s.idleConnsClosed
	log.Println("Server shutdown")
}

// HTTPServer returns the http server
func (s *Server) HTTPServer() *http.Server {
	return s.httpServer
}
