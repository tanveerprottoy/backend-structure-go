package e2e

import (
	"context"
	"log"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

// NewServer initializes the server
func NewServer(address string, handler http.Handler, opts ...any) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    address,
			Handler: handler,
		},
	}
}

func (s *Server) Shutdown() {
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP Server shutdown error: %v", err)
	}
	log.Println("Server shutdown")
}

// Start starts the server
func (s *Server) Start(chStart, chStop chan int) {
	// s.idleConnsClosed = make(chan struct{})
	log.Println("Server starting")
	// run in a separete goroutine
	go func() {
		// if err == http.ErrServerClosed do nothing
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Fatalf("HTTP Server ListenAndServe: %v", err)
		}
	}()
	// signal to the main routine that the server has started
	chStart <- 1
	// wait for stop signal
	<-chStop
	// received shutdown signal
	// commence shutdown
	s.Shutdown()
}

// HTTPServer returns the http server
func (s *Server) HTTPServer() *http.Server {
	return s.httpServer
}
