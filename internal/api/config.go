package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/delivery/http/handler"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/delivery/http/route"
	productcfg "github.com/tanveerprottoy/backend-structure-go/internal/api/product/config"
	usercfg "github.com/tanveerprottoy/backend-structure-go/internal/api/user/config"
	"github.com/tanveerprottoy/backend-structure-go/pkg/env"
	"github.com/tanveerprottoy/backend-structure-go/pkg/router"
	"github.com/tanveerprottoy/backend-structure-go/pkg/sqlext"
	"github.com/tanveerprottoy/backend-structure-go/pkg/validatorext"
)

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
	/* opts := sqlext.Options{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	} */
	// temp
	opts := sqlext.Options{
		Host:     "192.168.1.4",
		Port:     "5432",
		Username: "postgres",
		Password: "postgres",
		DBName:   "dummy_ecommerce_db",
		SSLMode:  "disable",
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

func (c *Config) initRoutes(handlers []any) {
	// handlers index
	// 0: product
	// 1: user
	productRoutes := route.Product(handlers[0].(*handler.Product))
	userRoutes := route.User(handlers[1].(*handler.User))

	// mount all the routes
	route.MountAll(
		c.router,
		[]any{
			productRoutes,
			userRoutes,
		},
	)
}

// initComponents initializes application components
func (c *Config) initComponents() {
	product := productcfg.NewConfig(c.DBClient.DB())
	user := usercfg.NewConfig(c.DBClient.DB())

	c.initRoutes(
		[]any{
			handler.NewProduct(product.UseCase, c.validater),
			handler.NewUser(user.UseCase, c.validater),
		},
	)
}
