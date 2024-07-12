package service_test

import (
	"context"
	"testing"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/user"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/user/mock"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/user/service"
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
		address := "address 1"

		tests := [2]struct {
			name     string
			dto      user.CreateDTO
			expected user.User
		}{
			{
				name: "success",
				dto: user.CreateDTO{
					Name:    "name 1",
					Address: "address 1",
				},
				expected: user.User{
					Name:    "name 1",
					Address: &address,
				},
			},
			{
				name: "fail",
				dto: user.CreateDTO{
					Name:    "name 2",
					Address: "address 2",
				},
				expected: user.User{
					Name:    "another name",
					Address: &address,
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
		address := "updated address 1"

		tests := [2]struct {
			name     string
			dto      user.UpdateDTO
			expected user.User
		}{
			{
				name: "success",
				dto: user.UpdateDTO{
					Name:    "updated name 1",
					Address: "updated address 1",
				},
				expected: user.User{
					ID:      insertedIDs[0],
					Name:    "updated name",
					Address: &address,
				},
			},
			{
				name: "fail",
				dto: user.UpdateDTO{
					Name:    "updated name 2",
					Address: "updated address 2",
				},
				expected: user.User{
					ID:      insertedIDs[1],
					Name:    "another name",
					Address: &address,
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
			expected user.User
		}{
			{
				name: "success",
				expected: user.User{
					ID: insertedIDs[0],
				},
			},
			{
				name: "fail",
				expected: user.User{
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
