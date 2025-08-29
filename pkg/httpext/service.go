package httpext

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type BasicAuth struct {
	UserName string
	Password string
}

func newRequest(
	ctx context.Context,
	method, url string,
	header http.Header,
	body io.Reader,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = header
	}

	return req, nil
}

func newRequestWithBasicAuth(
	ctx context.Context,
	method, url string,
	header http.Header,
	body io.Reader,
	basicAuth *BasicAuth,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = header
	}

	if basicAuth == nil {
		return nil, errors.New("basicAuth is nil")
	}
	// set basic auth
	req.SetBasicAuth(basicAuth.UserName, basicAuth.Password)

	return req, nil
}

// Request is a generic function to make a request with context
// Generic parameters: R = response type, E = error type
// use this function when you want to parse the response body to a specific type
// and also parse the error response to a specific type
// R and E must be a struct
func Request[R, E any](
	ctx context.Context,
	client Client,
	method string,
	url string,
	header http.Header,
	body io.Reader,
	retry bool,
	basicAuth *BasicAuth,
) (*R, *E, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), client.HTTPClient().Timeout)
		defer cancel()
	}

	var (
		req *http.Request
		err error
	)

	if basicAuth != nil {
		req, err = newRequestWithBasicAuth(ctx, method, url, header, body, basicAuth)
	} else {
		req, err = newRequest(ctx, method, url, header, body)
	}
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.Do(req, retry)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		// resp ok, parse response body to type
		var r R

		err := json.NewDecoder(resp.Body).Decode(&r)
		if err != nil {
			return nil, nil, err
		}

		return &r, nil, nil
	} else {
		// resp not ok, parse error
		var e E

		err := json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return nil, nil, err
		}

		return nil, &e, errors.New("error response was returned")
	}
}

// RequestRaw is a function to make a request with context
// it returns the status code and the response body as a byte slice
// use this function when need to get the raw response body
func RequestRaw(
	ctx context.Context,
	client Client,
	method string,
	url string,
	header http.Header,
	body io.Reader,
	retry bool,
	basicAuth *BasicAuth,
) (int, []byte, error) {
	var (
		req *http.Request
		err error
	)

	if basicAuth != nil {
		req, err = newRequestWithBasicAuth(ctx, method, url, header, body, basicAuth)
	} else {
		req, err = newRequest(ctx, method, url, header, body)
	}
	if err != nil {
		return 0, nil, err
	}

	resp, err := client.Do(req, retry)
	if err != nil {
		return 0, nil, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, b, nil
}
