package api

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/tanveerprottoy/backend-structure-go/pkg/env"
	"github.com/tanveerprottoy/backend-structure-go/pkg/router"
	"github.com/tanveerprottoy/backend-structure-go/pkg/sqlext"
	"github.com/tanveerprottoy/backend-structure-go/pkg/validatorext"
)

// Config contains the components of the application
// and configures them as required
type Config struct {
	DBClient  *sqlext.Client
	router    *router.Router
	validater validatorext.Validater
}

func NewConfig() *Config {
	c := new(Config)
	c.loadEnv()
	c.initDB()
	c.initRouter()
	c.initValidator()
	c.initComponents()
	return c
}

// loadEnv initializes env
func (c *Config) loadEnv() {
	env.LoadEnv("")
}

// initDB initializes DB client
func (c *Config) initDB() {
	opts := sqlext.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}

	c.DBClient = sqlext.GetInstance(opts)
}

// initRouter initializes router
func (c *Config) initRouter() {
	c.router = router.NewRouter()
}

// initValidator initializes validator
func (c *Config) initValidator() {
	c.validater = validatorext.NewValidator(validator.New())
}
