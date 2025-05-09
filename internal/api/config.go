package api

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/tanveerprottoy/backend-structure-go/pkg/env"
	"github.com/tanveerprottoy/backend-structure-go/pkg/router"
	"github.com/tanveerprottoy/backend-structure-go/pkg/sqlext"
	"github.com/tanveerprottoy/backend-structure-go/pkg/validatorext"
)

// config contains the components of the application
// and configures them as required
type config struct {
	dbClient  *sqlext.Client
	router    *router.Router
	validater validatorext.Validater
}

func NewConfig() *config {
	c := new(config)
	c.loadEnv()
	c.initDB()
	c.initRouter()
	c.initValidator()

	// init components
	initComponents(c)

	return c
}

// loadEnv initializes env
func (c *config) loadEnv() {
	env.LoadEnv("")
}

// initDB initializes DB client
func (c *config) initDB() {
	opts := sqlext.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}

	c.dbClient = sqlext.GetInstance(opts)
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
