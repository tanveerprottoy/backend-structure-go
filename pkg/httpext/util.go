package httpext

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetURLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func GetQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// ParseRequestBody parses the request body
func ParseRequestBody(r io.ReadCloser, v any) error {
	err := json.NewDecoder(r).Decode(&v)
	if err != nil {
		return err
	}
	return nil
}

func GetRouteContext(r *http.Request) *chi.Context {
	return chi.RouteContext(r.Context())
}
