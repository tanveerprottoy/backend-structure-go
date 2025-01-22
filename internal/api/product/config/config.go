package productcfg

import (
	"database/sql"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/product"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/product/postgres"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/product/service"
)

// Config contains and initializes the components of the product package
type Config struct {
	UseCase    product.UseCase
	Repository product.Repository
}

// NewConfig initializes a new NewConfig
func NewConfig(db *sql.DB) *Config {
	r := postgres.NewStorage(db)
	u := service.NewService(r)
	return &Config{UseCase: u, Repository: r}
}
