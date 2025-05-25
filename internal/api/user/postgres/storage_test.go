package postgres_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/user"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/user/postgres"
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
		addr := "Test User Address"

		n := time.Now().Unix()

		tests := [2]struct {
			name     string
			dto      *user.CreateDTO
			expected user.User
		}{
			{
				name: "success",
				dto: &user.CreateDTO{
					Name:      "Test User",
					Address:   &addr,
					CreatedAt: n,
					UpdatedAt: n,
				},
				expected: user.User{
					Name:    "Test User",
					Address: &addr,
				},
			},
			{
				name: "fail",
				dto: &user.CreateDTO{
					Name:    "name 2",
					Address: nil,
				},
				expected: user.User{
					Name:    "another name",
					Address: &addr,
				},
			},
		}

		for _, tc := range tests {
			mock.ExpectQuery(regexp.QuoteMeta(
				`INSERT INTO users (name, address, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`,
			)).
				WithArgs(tc.dto.Name, tc.dto.Address, tc.dto.CreatedAt, tc.dto.UpdatedAt).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

			gotID, err := s.Create(context.Background(), tc.dto)
			assert.NoError(t, err)

			assert.Equal(t, id, gotID)
		}
	})

	t.Run("ReadMany", func(t *testing.T) {
		// run test in parallel
		// t.Parallel()

		n := time.Now().Unix()

		// mock insert rows
		rows := sqlmock.NewRows([]string{"id", "name", "address", "is_archived", "created_at", "updated_at"}).
			AddRow(uuid.New().String(), "User1", "Address1", false, n, n).
			AddRow(uuid.New().String(), "User2", "Address2", false, n, n)

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
					`SELECT id, name, address, is_archived, created_at, updated_at FROM users LIMIT $1 OFFSET $2`,
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

		n := time.Now().Unix()

		// mock insert row
		row := sqlmock.NewRows([]string{"id", "name", "address", "is_archived", "created_at", "updated_at"}).
			AddRow(uuid.New().String(), "User1", "Address1", false, n, n)

		mock.ExpectQuery(regexp.QuoteMeta(
			"SELECT id, name, address, is_archived, created_at, updated_at FROM users WHERE id = $1 LIMIT $2",
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

		addr := "Updated User Address"

		// test cases
		tests := [2]struct {
			name     string
			dto      *user.UpdateDTO
			expected int64
		}{
			{
				name: "success update 1",
				dto: &user.UpdateDTO{
					Name:      "updated name 1",
					Address:   &addr,
					UpdatedAt: time.Now().Unix(),
				},
				// rows affected
				expected: 1,
				// user.User{
				// 	ID:          id,
				// 	Name:        "updated name",
				// 	Address: &addr,
				// 	UpdatedAt:   time.Now().Unix(),
				// },
			},
			{
				name: "success update 2",
				dto: &user.UpdateDTO{

					Name:      "updated name 2",
					Address:   &addr,
					UpdatedAt: time.Now().Unix(),
				},
				expected: 1,
			},
		}

		n := time.Now().Unix()

		// mock insert row
		_ = sqlmock.NewRows([]string{"id", "name", "address", "is_archived", "created_at", "updated_at"}).
			AddRow(uuid.New().String(), "User update 1", "Update Address 1", false, n, n)

		for i, tc := range tests {
			// run test in a sub test
			t.Run(tc.name, func(t *testing.T) {
				mock.ExpectExec(regexp.QuoteMeta(
					"UPDATE users SET name = $1, address = $2, updated_at = $3 WHERE id = $4",
				)).
					WithArgs(tc.dto.Name, tc.dto.Address, tc.dto.UpdatedAt, ids[i]).
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
				// user.User{
				// 	ID:          id,
				// 	Name:        "updated name",
				// 	Address: &addr,
				// 	UpdatedAt:   time.Now().Unix(),
				// },
			},
			{
				name:     "fail",
				id:       constant.FakeUUID,
				expected: 1,
			},
		}

		for _, tc := range tests {
			// run test in a sub test
			t.Run(tc.name, func(t *testing.T) {

				mock.ExpectExec(regexp.QuoteMeta(
					`UPDATE users SET is_archived = $1, updated_at = $2 WHERE id = $3`,
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

func TestStorage_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening stub db: %v", err)
	}
	defer db.Close()
	storage := postgres.NewStorage(db)

	addr := "dummy address"

	dto := &user.CreateDTO{
		Name:      "John",
		Address:   &addr,
		CreatedAt: 1234567890,
		UpdatedAt: 1234567890,
	}

	query := "INSERT INTO users"
	mock.ExpectQuery(query).
		WithArgs(dto.Name, dto.Address, dto.CreatedAt, dto.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))

	id, err := storage.Create(context.Background(), dto)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if id != "1" {
		t.Errorf("expected id '1', got '%s'", id)
	}
}
