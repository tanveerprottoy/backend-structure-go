package e2e

import (
	"github.com/tanveerprottoy/backend-structure-go/internal/api/delivery/http/handler"
	productcfg "github.com/tanveerprottoy/backend-structure-go/internal/api/product/config"
	usercfg "github.com/tanveerprottoy/backend-structure-go/internal/api/user/config"
)

// initComponents initializes application components
func initComponents(cfg *config) {
	product := productcfg.NewConfig(cfg.db)

	user := usercfg.NewConfig(cfg.db)

	initRoutes(
		cfg.router,
		[]any{
			handler.NewProduct(product.UseCase, cfg.validater),
			handler.NewUser(user.UseCase, cfg.validater),
		},
	)
}
