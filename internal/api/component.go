package api

import (
	"github.com/tanveerprottoy/backend-structure-go/internal/api/delivery/http/handler"
	productprovider "github.com/tanveerprottoy/backend-structure-go/internal/api/product/provider"
	userprovider "github.com/tanveerprottoy/backend-structure-go/internal/api/user/provider"
)

// initComponents initializes application components
func initComponents(cfg *config) {
	productProvider := productprovider.New(cfg.dbClient.DB())

	userProvider := userprovider.New(cfg.dbClient.DB())

	initRoutes(
		cfg.router,
		[]any{
			handler.NewProduct(productProvider.UseCase, cfg.validater),
			handler.NewUser(userProvider.UseCase, cfg.validater),
		},
	)
}
