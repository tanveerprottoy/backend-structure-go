package api

import (
	"github.com/tanveerprottoy/backend-structure-go/pkg/server"
)

// App contains the configuration and server
type App struct {
	cfg *Config
	srv *server.Server
}

// NewApp creates App
func NewApp() *App {
	a := &App{cfg: NewConfig()}
	a.initServer()
	a.configureGracefulShutdown()
	return a
}

// initServer initializes the server
func (a *App) initServer() {
	a.srv = server.NewServer(":8080", a.cfg.router.Mux)
	// a.srv = server.NewServer(":"+os.Getenv("PORT"), a.cfg.router.Mux)
}

// configureGracefulShutdown configures graceful shutdown
func (a *App) configureGracefulShutdown() {
	a.srv.ConfigureGracefulShutdown(func() {
		a.cfg.DBClient.Close()
	})
}

// Start starts the server
func (a *App) Start() {
	a.srv.Start()
}
