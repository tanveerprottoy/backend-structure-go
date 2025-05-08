package httpext

import (
	"net/http"
)

type Client interface {
	Do(req *http.Request, retry bool) (*http.Response, error)

	HTTPClient() *http.Client
}
