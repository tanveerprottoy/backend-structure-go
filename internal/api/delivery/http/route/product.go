package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/delivery/http/handler"
)

func Product(handler *handler.Product) chi.Router {
	r := chi.NewRouter()
	r.Post("/", handler.Create)
	r.Get("/", handler.ReadMany)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", handler.ReadOne)
		r.Put("/", handler.Update)
		r.Delete("/", handler.Delete)
	})
	return r
}
