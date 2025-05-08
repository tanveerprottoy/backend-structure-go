package e2e

import (
	"github.com/tanveerprottoy/backend-structure-go/internal/api/delivery/http/handler"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/delivery/http/route"
	"github.com/tanveerprottoy/backend-structure-go/pkg/router"
)

// initRoutes initializes all the routes
func initRoutes(router *router.Router, handlers []any) {
	// handlers index
	// 0: product
	// 1: user
	productRoutes := route.Product(handlers[0].(*handler.Product))
	userRoutes := route.User(handlers[1].(*handler.User))

	// mount all the routes
	route.MountAll(
		router,
		[]any{
			productRoutes,
			userRoutes,
		},
	)
}
