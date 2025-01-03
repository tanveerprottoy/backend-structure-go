package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/user"
	"github.com/tanveerprottoy/backend-structure-go/pkg/errorext"
	"github.com/tanveerprottoy/backend-structure-go/pkg/sqlext"
)

const tableName = "users"

// storage implements the storage interface
type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *storage {
	return &storage{db: db}
}

func (s *storage) Create(ctx context.Context, e user.User, args ...any) (string, error) {
	var lastID string

	// convert domain product entity to postgres specific entity
	entity := newUserEntity(e.Name, e.Address, e.CreatedAt, e.UpdatedAt)

	// build insert query
	q := sqlext.BuildInsertQuery(tableName, []string{"name", "address", "created_at", "updated_at"}, "RETURNING id")

	// execute the query
	row := s.db.QueryRowContext(ctx, q, entity.name, entity.address, entity.createdAt, entity.updatedAt)
	err := row.Err()
	if err != nil {
		log.Printf("err: %v", err)
		err := errorext.BuildDBError(err)
		return lastID, err
	}

	err = row.Scan(&lastID)
	if err != nil {
		log.Printf("err: %v", err)
		err := errorext.BuildDBError(err)
		return lastID, err
	}
	return lastID, nil
}

func (s *storage) ReadMany(ctx context.Context, limit, offset int, args ...any) ([]user.User, error) {
	cl := ""
	vals := make([]any, 0)
	d := make([]user.User, 0)
	if args[0] != nil {
		cl = " WHERE is_archived = $1"
		vals = append(vals, args[0].(bool))
	}
	if cl == "" {
		cl = " LIMIT $1 OFFSET $2"
	} else {
		cl += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(vals)+1, len(vals)+2)
	}
	vals = append(vals, limit, offset)
	q := "SELECT * FROM " + tableName + cl

	rows, err := s.db.QueryContext(ctx, q, vals...)
	if err != nil {
		err := errorext.BuildDBError(err)
		return d, err
	}

	defer rows.Close()

	// scan the rows
	entity := &userEntity{}
	users, err := entity.scanRows(rows)
	if err != nil {
		return d, err
	}

	// convert postgres entity to domain entity
	for _, u := range users {
		e := user.NewUser(u.id, u.name, u.address.String, u.createdAt, u.updatedAt)
		d = append(d, *e)
	}

	return d, nil
}

func (s *storage) ReadOne(ctx context.Context, id string, args ...any) (user.User, error) {
	q := sqlext.BuildSelectQuery(tableName, []string{}, []string{"id"}, "LIMIT $2")

	row := s.db.QueryRowContext(ctx, q, id, 1)
	err := row.Err()
	if err != nil {
		err := errorext.BuildDBError(err)
		return user.User{}, err
	}

	entity := &userEntity{}
	err = entity.scanRow(row)
	if err != nil {
		return user.User{}, err
	}

	// convert postgres entity to domain entity
	e := user.NewUser(entity.id, entity.name, entity.address.String, entity.createdAt, entity.updatedAt)

	return *e, nil
}

func (s *storage) Update(ctx context.Context, id string, e user.User, args ...any) (int64, error) {
	// convert domain product entity to postgres specific entity
	entity := newUserEntity(e.Name, e.Address, e.CreatedAt, e.UpdatedAt)

	q := sqlext.BuildUpdateQuery(tableName, []string{"name", "description", "updated_at"}, []string{"id"}, "")

	res, err := s.db.ExecContext(ctx, q, entity.name, entity.address, entity.updatedAt, id)
	if err != nil {
		err := errorext.BuildDBError(err)
		return -1, err
	}

	return sqlext.GetRowsAffected(res), nil
}

func (s *storage) Delete(ctx context.Context, id string, args ...any) (int64, error) {
	q := sqlext.BuildUpdateQuery(tableName, []string{"is_archived", "updated_at"}, []string{"id"}, "")

	res, err := s.db.ExecContext(ctx, q, true, args[0].(int64), id)
	if err != nil {
		err := errorext.BuildDBError(err)
		return -1, err
	}

	return sqlext.GetRowsAffected(res), nil
}
