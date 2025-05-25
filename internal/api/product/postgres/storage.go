package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/product"
	"github.com/tanveerprottoy/backend-structure-go/pkg/errorext"
	"github.com/tanveerprottoy/backend-structure-go/pkg/sqlext"
)

const tableName = "products"

// storage implements the storage interface
type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *storage {
	return &storage{db: db}
}

func (s *storage) Create(ctx context.Context, dto *product.CreateDTO, args ...any) (string, error) {
	var lastID string

	// build insert query
	q := sqlext.BuildInsertQuery(tableName, []string{"name", "description", "created_at", "updated_at"}, "RETURNING id")

	// execute the query
	row := s.db.QueryRowContext(ctx, q, dto.Name, dto.Description, dto.CreatedAt, dto.UpdatedAt)
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

func (s *storage) ReadMany(ctx context.Context, limit, offset int, args ...any) ([]product.Product, error) {
	d := make([]product.Product, 0)

	q := fmt.Sprintf("SELECT id, name, description, is_archived, created_at, updated_at FROM %s", tableName)
	vals := make([]any, 0)

	if len(args) > 0 && args[0] != nil {
		q += " WHERE is_archived = $1"
		vals = append(vals, args[0].(bool))
	}

	q += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(vals)+1, len(vals)+2)
	vals = append(vals, limit, offset)

	rows, err := s.db.QueryContext(ctx, q, vals...)
	if err != nil {
		err := errorext.BuildDBError(err)
		return d, err
	}

	defer rows.Close()

	// scan the rows
	entity := &productEntity{}
	products, err := entity.scanRows(rows)
	if err != nil {
		return d, err
	}

	// convert postgres entity to domain entity
	for _, p := range products {
		d = append(d, product.Product{
			ID:          p.id,
			Name:        p.name,
			Description: &p.description.String,
			IsArchived:  p.isArchived,
			CreatedAt:   p.createdAt,
			UpdatedAt:   p.updatedAt,
		})
	}

	return d, nil
}

func (s *storage) ReadOne(ctx context.Context, id string, args ...any) (product.Product, error) {
	projections := []string{"id", "name", "description", "is_archived", "created_at", "updated_at"}

	q := sqlext.BuildSelectQuery(tableName, projections, []string{"id"}, "LIMIT $2")

	row := s.db.QueryRowContext(ctx, q, id, 1)
	err := row.Err()
	if err != nil {
		err := errorext.BuildDBError(err)
		return product.Product{}, err
	}

	entity := &productEntity{}
	err = entity.scanRow(row)
	if err != nil {
		return product.Product{}, err
	}

	// convert postgres entity to domain entity
	return product.Product{
		ID:          entity.id,
		Name:        entity.name,
		Description: &entity.description.String,
		IsArchived:  entity.isArchived,
		CreatedAt:   entity.createdAt,
		UpdatedAt:   entity.updatedAt,
	}, nil
}

func (s *storage) Update(ctx context.Context, id string, dto *product.UpdateDTO, args ...any) (int64, error) {
	q := sqlext.BuildUpdateQuery(tableName, []string{"name", "description", "updated_at"}, []string{"id"}, "")

	res, err := s.db.ExecContext(ctx, q, dto.Name, dto.Description, dto.UpdatedAt, id)
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
