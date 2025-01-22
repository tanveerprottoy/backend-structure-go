package service_test

import (
	"context"
	"testing"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/product"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/product/mock"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/product/service"
	"github.com/tanveerprottoy/backend-structure-go/pkg/constant"
)

func TestService(t *testing.T) {
	// create mock repo
	r := mock.NewMemoryStorage()

	s := service.NewService(r)

	// cleanup
	t.Cleanup(func() {
		r.Clear()
	})

	var insertedIDs [2]string

	// initiate the tests in sub tests
	t.Run("create", func(t *testing.T) {
		descrption := "description 1"

		tests := [2]struct {
			name     string
			dto      *product.CreateDTO
			expected product.Product
		}{
			{
				name: "success",
				dto: &product.CreateDTO{
					Name:        "name 1",
					Description: nil,
				},
				expected: product.Product{
					Name:        "name 1",
					Description: &descrption,
				},
			},
			{
				name: "fail",
				dto: &product.CreateDTO{
					Name:        "name 2",
					Description: nil,
				},
				expected: product.Product{
					Name:        "another name",
					Description: &descrption,
				},
			},
		}

		for i, test := range tests {
			// run test in a sub test
			t.Run(test.name, func(t *testing.T) {
				e, err := s.Create(context.Background(), test.dto)
				if err != nil {
					t.Error(err)
				}
				if e.Name != test.expected.Name {
					t.Error("name does not match")
				}
				insertedIDs[i] = e.ID
			})
		}
	})

	t.Run("readMany", func(t *testing.T) {
		if len(insertedIDs) == 0 {
			t.Skip("no inserted id found, skipping test")
		}

		tests := [2]struct {
			name     string
			expected int
		}{
			{
				name:     "success",
				expected: 2, // expected item count 2
			},
			{
				name:     "fail",
				expected: 3, // expected item count 3
			},
		}

		for _, test := range tests {
			// run test in a sub test
			t.Run(test.name, func(t *testing.T) {
				d, err := s.ReadMany(context.Background(), 10, 1)
				if err != nil {
					t.Error(err)
				}
				l := len(d)
				if l == 0 {
					t.Error("no date returned")
				}
				if l != test.expected {
					t.Errorf("expected item count %d, got %d", l, test.expected)
				}
			})
		}
	})

	t.Run("readOne", func(t *testing.T) {
		if len(insertedIDs) == 0 {
			t.Skip("no inserted id found, skipping test")
		}

		id := insertedIDs[0]

		tests := [2]struct {
			name     string
			expected string
		}{
			{
				name:     "success",
				expected: id,
			},
			{
				name:     "fail",
				expected: constant.FakeUUID,
			},
		}

		for _, test := range tests {
			// run test in a sub test
			t.Run(test.name, func(t *testing.T) {
				e, err := s.ReadOne(context.Background(), id)
				if err != nil {
					t.Error(err)
				}
				if e.ID != test.expected {
					t.Error("id is not equal")
				}
			})
		}
	})

	t.Run("update", func(t *testing.T) {
		descrption := "updated description 1"

		tests := [2]struct {
			name     string
			dto      *product.UpdateDTO
			expected product.Product
		}{
			{
				name: "success",
				dto: &product.UpdateDTO{
					Name:        "updated name 1",
					Description: nil,
				},
				expected: product.Product{
					ID:          insertedIDs[0],
					Name:        "updated name",
					Description: &descrption,
				},
			},
			{
				name: "fail",
				dto: &product.UpdateDTO{
					Name:        "updated name 2",
					Description: nil,
				},
				expected: product.Product{
					ID:          insertedIDs[1],
					Name:        "another name",
					Description: &descrption,
				},
			},
		}

		// simulate failure by supplying wrong id
		ids := [2]string{insertedIDs[0], constant.FakeUUID}

		for i, test := range tests {
			// run test in a sub test
			t.Run(test.name, func(t *testing.T) {
				e, err := s.Update(context.Background(), ids[i], test.dto)
				if err != nil {
					t.Error(err)
				}
				if e.ID != test.expected.ID {
					t.Error("id does not match")
				}
			})
		}
	})

	t.Run("delete", func(t *testing.T) {
		// simulate failure by supplying wrong id
		ids := [2]string{insertedIDs[0], constant.FakeUUID}

		tests := [2]struct {
			name     string
			expected product.Product
		}{
			{
				name: "success",
				expected: product.Product{
					ID: insertedIDs[0],
				},
			},
			{
				name: "fail",
				expected: product.Product{
					ID: insertedIDs[1],
				},
			},
		}

		for i, test := range tests {
			// run test in a sub test
			t.Run(test.name, func(t *testing.T) {
				e, err := s.Delete(context.Background(), ids[i])
				if err != nil {
					t.Error(err)
				}
				if e.ID != test.expected.ID {
					t.Error("id does not match")
				}
			})
		}
	})
}
