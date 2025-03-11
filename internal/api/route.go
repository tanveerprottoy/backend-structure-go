package api

import (
	"github.com/tanveerprottoy/backend-structure-go/internal/api/delivery/http/handler"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/delivery/http/route"
)

// initRoutes initializes all the routes
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
