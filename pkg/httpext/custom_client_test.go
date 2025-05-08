package httpext_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/tanveerprottoy/backend-structure-go/pkg/httpext"
)

func TestCustomClient(t *testing.T) {
	cfg := httpext.Config{
		MaxRetries: 5,
		MaxJitter:  10,
		Timeout:    10 * time.Second,
	}

	client := httpext.NewCustomClient(
		cfg,
		httpext.WithIdleConnTimeout(50*time.Second),
		httpext.WithMaxIdleConnsPerHost(20),
	)

	// subtest testDo
	t.Run("testDo", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		// Mock a request
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/api/v1/products", nil)
		if err != nil {
			t.Errorf("failed to create request: %v", err)
			return
		}

		resp, err := client.Do(req, true)
		if err != nil {
			t.Errorf("client.Do error: %v", err)
			return
		}

		if resp == nil {
			t.Errorf("client.Do returned nil response")
			return
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code 200, got %d", resp.StatusCode)
			return
		}
	})

}
