package e2e

import (
	"github.com/tanveerprottoy/backend-structure-go/internal/api/delivery/http/handler"
	productprovider "github.com/tanveerprottoy/backend-structure-go/internal/api/product/provider"
	userprovider "github.com/tanveerprottoy/backend-structure-go/internal/api/user/provider"
)

// initComponents initializes application components
func initComponents(cfg *config) {
	productProvider := productprovider.New(cfg.db)

	userProvider := userprovider.New(cfg.db)

	initRoutes(
		cfg.router,
		[]any{
			handler.NewProduct(productProvider.UseCase, cfg.validater),
			handler.NewUser(userProvider.UseCase, cfg.validater),
		},
	)
}
