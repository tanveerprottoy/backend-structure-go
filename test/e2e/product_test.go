package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/product"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/user"
	"github.com/tanveerprottoy/backend-structure-go/pkg/constant"
	"github.com/tanveerprottoy/backend-structure-go/pkg/httpext"
	"github.com/tanveerprottoy/backend-structure-go/pkg/response"
)

func TestProduct(t *testing.T) {
	ctx := context.Background()

	desc := "test description"

	// Mock data
	n := time.Now().Unix()

	e := product.Product{
		Name:        "Test",
		Description: &desc,
		CreatedAt:   n,
		UpdatedAt:   n,
	}

	// build test table data
	/* var tests = []struct {
		name string
		data user.User
		expected user.User
		got user.User
	}{
		{
			name: "Create",
			data:e,
			expected:e,
		},
	} */
	// test create
	t.Run("create", func(t *testing.T) {
		b, err := json.Marshal(e)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		u, errRes, err := httpext.Request[response.Response[user.User], response.ErrorResponse](ctx, httpClient, http.MethodPost, baseURL+constant.V1+constant.UsersPattern, nil, bytes.NewReader(b), false)
		// log.Printf("create res: %v\n", res)
		if err != nil {
			// check if errRes has error
			if errRes != nil {
				log.Println("errRes: ", errRes)
			}

			t.Errorf("error returned: %v", err)
		}

		e.ID = u.Data.ID
	})

	t.Run(("readMany"), func(t *testing.T) {
		resp, errRes, err := httpext.Request[response.Response[response.ReadManyResponse[user.User]], response.ErrorResponse](ctx, httpClient, http.MethodGet, baseURL+constant.V1+constant.UsersPattern, nil, nil, false)
		// log.Printf("readMany res: %v\n", res)
		if err != nil {
			// check if errRes has error
			if errRes != nil {
				log.Println("errRes: ", errRes)
			}

			t.Errorf("error returned: %v", err)
		}

		if len(resp.Data.Items) == 0 {
			t.Error("empty data returned, expected one item")
		}
	})

	t.Run(("readOne"), func(t *testing.T) {
		resp, errRes, err := httpext.Request[response.Response[user.User], response.ErrorResponse](ctx, httpClient, http.MethodGet, baseURL+constant.V1+constant.UsersPattern+"/"+e.ID, nil, nil, false)
		// log.Printf("readMany res: %v\n", res)
		if err != nil {
			// check if errRes has error
			if errRes != nil {
				log.Println("errRes: ", errRes)
			}

			t.Errorf("error returned: %v", err)
		}

		if resp.Data.ID != e.ID {
			t.Errorf("expected id %s, got %s", e.ID, resp.Data.ID)
		}
	})

	t.Run("update", func(t *testing.T) {
		e.Name = "test update"
		b, err := json.Marshal(e)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		resp, errRes, err := httpext.Request[response.Response[user.User], response.ErrorResponse](ctx, httpClient, http.MethodPut, baseURL+constant.V1+constant.UsersPattern+"/"+e.ID, nil, bytes.NewReader(b), false)
		// log.Printf("readMany res: %v\n", res)
		if err != nil {
			// check if errRes has error
			if errRes != nil {
				log.Println("errRes: ", errRes)
			}

			t.Errorf("error returned: %v", err)
		}

		if resp.Data.Name != e.Name {
			t.Errorf("expected value %s, got %s", e.Name, resp.Data.Name)
		}
	})

	t.Run("delete", func(t *testing.T) {
		resp, errRes, err := httpext.Request[response.Response[user.User], response.ErrorResponse](ctx, httpClient, http.MethodDelete, baseURL+constant.V1+constant.UsersPattern+"/"+e.ID, nil, nil, false)
		// log.Printf("readMany res: %v\n", res)
		if err != nil {
			// check if errRes has error
			if errRes != nil {
				log.Println("errRes: ", errRes)
			}

			t.Errorf("error returned: %v", err)
		}

		if !resp.Data.IsArchived {
			t.Errorf("expected IsArchived true, got false")
		}
	})
}
