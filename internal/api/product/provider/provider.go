package provider

import (
	"database/sql"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/product"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/product/postgres"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/product/service"
)

// Provider contains and initializes the components of the package
type Provider struct {
	UseCase    product.UseCase
	Repository product.Repository
}

// New initializes a Provider
func New(db *sql.DB) Provider {
	r := postgres.NewStorage(db)
	u := service.NewService(r)
	return Provider{UseCase: u, Repository: r}
}
