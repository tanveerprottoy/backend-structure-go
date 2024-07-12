package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/tanveerprottoy/backend-structure-go/pkg/constant"
	"github.com/tanveerprottoy/backend-structure-go/pkg/router"
)

// this will contain all the routes for the application
// MountAll will mount all the routes
// their respective handlers will be passed as an argument
func MountAll(router *router.Router, routes []any) {
	// routes index
	// 0: product
	// 1: user
	router.Mux.Mount(constant.ApiPattern, router.Mux.Group(
		func(r chi.Router) {
			// v1 routes
			// 0 contains product routes
			r.Mount(constant.V1+constant.ProductsPattern, routes[0].(chi.Router))
			// 1 contains user routes
			r.Mount(constant.V1+constant.UsersPattern, routes[1].(chi.Router))
		}),
	)
}
