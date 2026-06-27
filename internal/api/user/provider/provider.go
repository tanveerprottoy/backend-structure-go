package provider

import (
	"database/sql"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/user"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/user/postgres"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/user/service"
)

// Provider contains and initializes the components of the package
type Provider struct {
	UseCase    user.UseCase
	Repository user.Repository
}

// New initializes a Provider
func New(db *sql.DB) Provider {
	r := postgres.NewStorage(db)
	u := service.NewService(r)
	return Provider{UseCase: u, Repository: r}
}
