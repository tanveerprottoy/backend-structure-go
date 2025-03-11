package api

import (
	"github.com/tanveerprottoy/backend-structure-go/internal/api/delivery/http/handler"
	productcfg "github.com/tanveerprottoy/backend-structure-go/internal/api/product/config"
	usercfg "github.com/tanveerprottoy/backend-structure-go/internal/api/user/config"
)

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
