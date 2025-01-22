package usercfg

import (
	"database/sql"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/user"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/user/postgres"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/user/service"
)

// Config contains and initializes the components of the user package
type Config struct {
	UseCase    user.UseCase
	Repository user.Repository
}

// NewConfig initializes a new NewConfig
func NewConfig(db *sql.DB) *Config {
	r := postgres.NewStorage(db)
	u := service.NewService(r)
	return &Config{UseCase: u, Repository: r}
}
