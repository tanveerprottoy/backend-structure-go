package service

import (
	"context"
	"testing"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/product"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/product/mock"
	"github.com/tanveerprottoy/backend-structure-go/pkg/constant"
)

func TestService(t *testing.T) {
	// create mock repo
	r := mock.NewMemoryStorage()

	s := NewService(r)

	// cleanup
	t.Cleanup(func() {
		r.Clear()
	})

	// initiate the tests in sub tests
	t.Run("readOneInternal", func(t *testing.T) {
		var insertedIDs [2]string

		address := "address 1"

		dtos := []*product.CreateDTO{
			{
				Name:        "name 1",
				Description: nil,
			},
			{
				Name:        "name 2",
				Description: &address,
			},
		}

		for i, dto := range dtos {
			// insert item for test
			e, err := s.Create(context.Background(), dto)
			if err != nil {
				t.Skip(err)
			}

			insertedIDs[i] = e.ID
		}

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
				e, err := s.readOneInternal(context.Background(), id)
				if err != nil {
					t.Error(err)
				}
				
				if e.ID != test.expected {
					t.Error("id is not equal")
				}
			})
		}
	})
}
