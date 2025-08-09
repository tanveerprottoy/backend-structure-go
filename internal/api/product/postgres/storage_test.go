package postgres_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/product"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/product/postgres"
	"github.com/tanveerprottoy/backend-structure-go/pkg/constant"
)

func TestStorage(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	t.Cleanup(func() {
		defer db.Close()
	})

	s := postgres.NewStorage(db)

	// inserted ids stored for later use
	// var insertedIDs [2]string

	t.Run("Create", func(t *testing.T) {
		// run test in parallel
		// t.Parallel()

		id := uuid.New().String()
		desc := "Test Product Description"

		n := time.Now().Unix()

		tests := [2]struct {
			name     string
			dto      *product.CreateDTO
			expected product.Product
		}{
			{
				name: "success 1",
				dto: &product.CreateDTO{
					Name:        "Test Product",
					Description: &desc,
					CreatedAt:   n,
					UpdatedAt:   n,
				},
				expected: product.Product{
					Name:        "Test Product",
					Description: &desc,
				},
			},
			{
				name: "success 2",
				dto: &product.CreateDTO{
					Name:        "name 2",
					Description: nil,
				},
				expected: product.Product{
					Name:        "name 3",
					Description: &desc,
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				mock.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO products (name, description, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`,
				)).
					WithArgs(tc.dto.Name, tc.dto.Description, tc.dto.CreatedAt, tc.dto.UpdatedAt).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

				gotID, err := s.Create(context.Background(), tc.dto)
				assert.NoError(t, err)
				assert.Equal(t, id, gotID)
			})
		}
	})

	t.Run("ReadMany", func(t *testing.T) {
		// run test in parallel
		// t.Parallel()

		n := time.Now().Unix()

		// mock insert rows
		rows := sqlmock.NewRows([]string{"id", "name", "description", "is_archived", "created_at", "updated_at"}).
			AddRow(uuid.New().String(), "Product1", "description 1", false, n, n).
			AddRow(uuid.New().String(), "Product2", "description 2", false, n, n)

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

		for _, tc := range tests {
			// run test in a sub test
			t.Run(tc.name, func(t *testing.T) {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT id, name, description, is_archived, created_at, updated_at FROM products LIMIT $1 OFFSET $2",
				)).
					WithArgs(2, 0).
					WillReturnRows(rows)

				d, err := s.ReadMany(context.Background(), 2, 0)
				assert.NoError(t, err)
				assert.Len(t, d, tc.expected)
			})
		}
	})

	t.Run("ReadOne", func(t *testing.T) {
		// run test in parallel
		// t.Parallel()

		id := uuid.New().String()

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

		// mock insert row
		row := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
			AddRow(id, "Product1", "description 1", time.Now().Unix(), time.Now().Unix())

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, name, description, created_at, updated_at FROM products WHERE id = $1 LIMIT $2`,
		)).
			WithArgs(id, 1).
			WillReturnRows(row)

		for _, tc := range tests {
			// run test in a sub test
			t.Run(tc.name, func(t *testing.T) {
				e, err := s.ReadOne(context.Background(), id)
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, e.ID)
			})
		}
	})

	t.Run("Update", func(t *testing.T) {
		// run test in parallel
		// t.Parallel()

		id := uuid.New().String()

		// mock ids
		ids := [2]string{
			id,
			constant.FakeUUID,
		}

		desc := "Updated Product Description"

		// test cases
		tests := [2]struct {
			name     string
			dto      *product.UpdateDTO
			expected int64
		}{
			{
				name: "success",
				dto: &product.UpdateDTO{
					Name:        "updated name 1",
					Description: &desc,
					UpdatedAt:   time.Now().Unix(),
				},
				// rows affected
				expected: 1,
				// product.Product{
				// 	ID:          id,
				// 	Name:        "updated name",
				// 	Description: &desc,
				// 	UpdatedAt:   time.Now().Unix(),
				// },
			},
			{
				name: "fail",
				dto: &product.UpdateDTO{
					Name:        "updated name 2",
					Description: &desc,
					UpdatedAt:   time.Now().Unix(),
				},
				expected: 0,
			},
		}

		for i, tc := range tests {
			// run test in a sub test
			t.Run(tc.name, func(t *testing.T) {
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE products SET name = $1, description = $2, updated_at = $3 WHERE id = $4`,
				)).
					WithArgs(tc.dto.Name, tc.dto.Description, tc.dto.UpdatedAt, ids[i]).
					WillReturnResult(sqlmock.NewResult(0, 1))

				rowsAffected, err := s.Update(context.Background(), ids[i], tc.dto)
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, rowsAffected)
			})
		}
	})

	t.Run("Delete", func(t *testing.T) {
		// run test in parallel
		// t.Parallel()

		// test cases
		tests := [2]struct {
			name     string
			id       string
			expected int64
		}{
			{
				name: "success",
				id:   uuid.New().String(),
				// rows affected
				expected: 1,
				// product.Product{
				// 	ID:          id,
				// 	Name:        "updated name",
				// 	Description: &desc,
				// 	UpdatedAt:   time.Now().Unix(),
				// },
			},
			{
				name:     "fail",
				id:       constant.FakeUUID,
				expected: 0,
			},
		}

		for _, tc := range tests {
			// run test in a sub test
			t.Run(tc.name, func(t *testing.T) {
				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE products SET is_archived = $1, updated_at = $2 WHERE id = $3`,
				)).
					WithArgs(true, sqlmock.AnyArg(), tc.id).
					WillReturnResult(sqlmock.NewResult(0, 1))

				rowsAffected, err := s.Delete(context.Background(), tc.id, time.Now().Unix())
				assert.NoError(t, err)
				assert.Equal(t, int64(1), rowsAffected)
			})
		}
	})
}
