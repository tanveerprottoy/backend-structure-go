package router

import (
	"net/http"
	"time"

	"github.com/tanveerprottoy/backend-structure-go/pkg/constant"
	middlewarext "github.com/tanveerprottoy/backend-structure-go/pkg/httpext/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Router struct
type Router struct {
	Mux *chi.Mux
}

func NewRouter() *Router {
	r := &Router{}
	r.Mux = chi.NewRouter()
	r.registerGlobalMiddlewares()
	r.registerNotFoundHander()
	return r
}

func (r *Router) registerGlobalMiddlewares() {
	r.Mux.Use(
		middleware.Logger,
		middleware.Recoverer,
		middlewarext.JSONContentTypeMiddleWare,
		// timeout middlewares
		middleware.Timeout(constant.RequestTimeout*time.Second),
		middlewarext.TimeoutHandler(constant.RequestTimeout*time.Second),
		cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"*"},
			// AllowedOrigins: []string{os.Getenv("ALLOWED_ORIGIN")},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods: constant.AllowedMethods,
			// AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowedHeaders: constant.AllowedHeaders,
			// ExposedHeaders:   []string{"Link"},
			AllowCredentials: constant.AllowCredentials,
			// MaxAge:           300, // Maximum value not ignored by any of major browsers
		}),
		// middlewarepkg.CORSEnableMiddleWareChi,
		// middlewarepkg.CORSEnableMiddleWare,
	)
}

func (r *Router) registerNotFoundHander() {
	// not found handler
	r.Mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("{\"message\": \"Resource not found\"}"))
	})
}
