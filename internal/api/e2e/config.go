package e2e

import (
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/tanveerprottoy/backend-structure-go/pkg/router"
	"github.com/tanveerprottoy/backend-structure-go/pkg/validatorext"
)

// config contains the components of the application
// and configures them as required
type config struct {
	db        *sql.DB
	router    *router.Router
	validater validatorext.Validater
}

func NewConfig(db *sql.DB) *config {
	c := &config{db: db}
	c.initRouter()
	c.initValidator()

	// Initialize components
	initComponents(c)

	return c
}

// initRouter initializes router
func (c *config) initRouter() {
	c.router = router.NewRouter()
}

// initValidator initializes validator
func (c *config) initValidator() {
	c.validater = validatorext.NewValidator(validator.New())
}

func (c *config) Router() *router.Router {
	return c.router
}
