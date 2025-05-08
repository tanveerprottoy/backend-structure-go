package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/delivery/http/handler"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/user"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/user/postgres"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/user/service"
	"github.com/tanveerprottoy/backend-structure-go/pkg/response"
	"github.com/tanveerprottoy/backend-structure-go/pkg/timeext"
)

func TestUser(t *testing.T) {
	// init storage
	r := postgres.NewStorage(db)
	// init service
	s := service.NewService(r)
	// init handler
	h := handler.NewUser(s, validater)
	// Mock data
	n := timeext.NowUnix()

	addr := "test address"

	e := user.User{
		Name:      "Test",
		Address:   &addr,
		CreatedAt: n,
		UpdatedAt: n,
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
		b, _ := json.Marshal(e)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))

		// call
		h.Create(w, r)
		// evaluate response
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, res.StatusCode)
		}

		var resBody response.Response[user.User]
		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			t.Errorf("decode response body: %v", err)
		}

		e.ID = resBody.Data.ID

		if !reflect.DeepEqual(e, resBody.Data) {
			t.Errorf("got { %+v } did not match expected { %+v }", resBody.Data, e)
		}
	})

	t.Run(("readMany"), func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/?limit=10&page=1", nil)

		// call
		h.ReadMany(w, r)
		// evaluate response
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, res.StatusCode)
		}

		var resBody response.Response[response.ReadManyResponse[user.User]]
		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			t.Errorf("decode response body: %v", err)
		}
		if len(resBody.Data.Items) == 0 {
			t.Errorf("expected data, got empty")
		}
		if resBody.Data.Items[0].ID != e.ID {
			t.Errorf("expected id %s, got %s", e.ID, resBody.Data.Items[0].ID)
		}
	})

	t.Run(("readOne"), func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/{id}", nil)
		// add path param
		r = addPathParam(r, "id", e.ID)

		// call
		h.ReadOne(w, r)
		// evaluate response
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, res.StatusCode)
		}

		var resBody response.Response[user.User]
		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			t.Errorf("decode response body: %v", err)
		}
		if resBody.Data.ID != e.ID {
			t.Errorf("expected id %s, got %s", e.ID, resBody.Data.ID)
		}
	})

	t.Run("update", func(t *testing.T) {
		e.Name = "Test Update"
		b, _ := json.Marshal(e)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/{id}", bytes.NewReader(b))
		// add path param
		r = addPathParam(r, "id", e.ID)

		// call
		h.Update(w, r)
		// evaluate response
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, res.StatusCode)
		}

		var resBody response.Response[user.User]
		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			t.Errorf("decode response body: %v", err)
		}
		if !reflect.DeepEqual(e, resBody.Data) {
			t.Errorf("got { %+v } did not match expected { %+v }", resBody.Data, e)
		}
	})

	t.Run("delete", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/{id}", nil)
		// add path param
		r = addPathParam(r, "id", e.ID)

		// call
		h.Delete(w, r)
		// evaluate response
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, res.StatusCode)
		}

		var resBody response.Response[user.User]
		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			t.Errorf("decode response body: %v", err)
		}
		if !resBody.Data.IsArchived {
			t.Errorf("expected IsArchived true, got false")
		}
	})
}
