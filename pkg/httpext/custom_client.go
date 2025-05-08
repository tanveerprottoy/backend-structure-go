package httpext

import (
	"bytes"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

type Config struct {
	MaxRetries int           // maximum number of retries for a request
	MaxJitter  int           // maximum jitter in milliseconds
	Timeout    time.Duration // request timeout
}

type Option func(*customClient)

func WithCheckRedirectFunc(f func(req *http.Request, via []*http.Request) error) Option {
	return func(c *customClient) {
		c.httpClient.CheckRedirect = f
	}
}

func WithMaxIdleConnsPerHost(maxIdleConnsPerHost int) Option {
	return func(c *customClient) {
		c.maxIdleConnsPerHost = maxIdleConnsPerHost
	}
}

func WithIdleConnTimeout(idleConnTimeout time.Duration) Option {
	return func(c *customClient) {
		c.idleConnTimeout = idleConnTimeout
	}
}

// customClient is a custom HTTP client that implements the Client interface
type customClient struct {
	httpClient *http.Client
	maxRetries int
	maxJitter  int

	// transport options
	maxIdleConnsPerHost int
	idleConnTimeout     time.Duration
}

func NewCustomClient(cfg Config, opts ...Option) *customClient {
	// sanitize maxRetries, maxJitter and maxIdleConnsPerHost
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 3
	}

	if cfg.MaxJitter <= 0 {
		cfg.MaxJitter = 10
	}

	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}

	httpClient := &http.Client{Timeout: cfg.Timeout}

	c := &customClient{
		httpClient: httpClient,
		maxRetries: cfg.MaxRetries,
		maxJitter:  cfg.MaxJitter,
	}

	// apply options
	for _, opt := range opts {
		opt(c)
	}

	// if one of the transport options is set, use the custom transport/roundtripper
	if c.maxIdleConnsPerHost > 0 || c.idleConnTimeout > 0 {
		httpClient.Transport = &http.Transport{
			MaxIdleConnsPerHost: c.maxIdleConnsPerHost,
			IdleConnTimeout:     c.idleConnTimeout,
		}
	}

	return c
}

// backoff generates the backoff time in seconds based on the number of retries
func (c *customClient) backoff(retries int) time.Duration {
	// 2^n backoff, n = number of retries
	return time.Duration(math.Pow(2, float64(retries))) * time.Second
}

// jitter generates a random jitter in milliseconds which is added to the backoff time
func (c *customClient) jitter(max, attempts int) time.Duration {
	rnd := rand.Intn(max)

	return time.Duration(int(attempts*rnd)) * time.Millisecond
}

func (c *customClient) restoreRequestBody(req *http.Request) error {
	if req.Body == nil {
		return nil
	}

	// Read the body into a buffer
	buf, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	// Restore the body so it can be read again
	// This is important because the body is an io.ReadCloser and can only be read once
	req.Body = io.NopCloser(bytes.NewBuffer(buf))

	return nil
}

func (c *customClient) isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// check if error is temporary
	if errNet, ok := err.(interface{ Temporary() bool }); ok && errNet.Temporary() {
		return true
	}

	return false
}

func (c *customClient) drainBody(resp *http.Response) {
	// drain the response body to reuse the connection
	// only do this if the response is not nil and the body is not nil
	if resp != nil && resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func (c *customClient) doWithRetry(req *http.Request) (*http.Response, error) {
	var (
		attempts = 0
		resp     *http.Response
		err      error
	)

	for attempts < c.maxRetries {
		log.Printf("customClient: doWithRetry called, attempt %d\n", attempts)

		// reusing a request body can be a bit tricky because the
		// io.ReadCloser interface, which is the type of r.Body in an
		// http.Request, is designed for single consumption. Once you've read the body, the underlying reader is often at its end, and attempting to read it again will yield an empty result or an error.
		// Always rewind/restore the request body when non-nil.
		if req.Body != nil {
			if err := c.restoreRequestBody(req); err != nil {
				return nil, err
			}
		}

		resp, err = c.httpClient.Do(req)

		log.Printf("customClient: doWithRetry.httpClient.Do called, attempt %d, response: %v\nerr: %v\n", attempts, resp, err)

		if err != nil {
			log.Printf("customClient: doWithRetry.httpClient.Do if err != nil: %v\n", err)

			// check if error is temporary
			if c.isRetryableError(err) {
				// temp
				log.Printf("backoff: %v\njitter: %v\n", c.backoff(attempts), c.jitter(c.maxJitter, attempts))

				// drain the response body to reuse the connection
				c.drainBody(resp)

				// wait for backoff time
				time.Sleep(c.backoff(attempts) + c.jitter(c.maxJitter, attempts))

				// continue to next attempt
				continue
			}

			return nil, err
		}

		// increment attempts
		attempts++

		log.Printf("customClient: doWithRetry after increment attempts: %d\n", attempts)
	}

	return resp, nil
}

func (c *customClient) doWithoutRetry(req *http.Request) (*http.Response, error) {
	// do without retry
	log.Println("customClient: doWithoutRetry called")

	return c.httpClient.Do(req)
}

func (c *customClient) Do(req *http.Request, retry bool) (*http.Response, error) {
	log.Println("customClient: Do called")
	if retry {
		return c.doWithRetry(req)
	}

	// do without retry
	return c.doWithoutRetry(req)

}

func (c *customClient) HTTPClient() *http.Client {
	return c.httpClient
}
