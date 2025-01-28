package postgres

import (
	"database/sql"
	"log"

	"github.com/tanveerprottoy/backend-structure-go/pkg/errorext"
)

// a package private entity clone of the domain entity
// this type can have db specific data types like here
// it's sqlext.NullString
type productEntity struct {
	id          string         `db:"id"`
	name        string         `db:"name"`
	description sql.NullString `db:"description"`
	isArchived  bool           `db:"is_archived"`
	createdAt   int64          `db:"created_at"`
	updatedAt   int64          `db:"updated_at"`
}

// this will be used to create db entity from domain entity
func newProductEntity(name string, description *string, createdAt, updatedAt int64) *productEntity {
	// description can be nil in db
	// check for nil
	e := productEntity{
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}

	if description != nil {
		e.description = sql.NullString{String: *description, Valid: true}
	}

	return &e
}

func (e *productEntity) scanRow(row *sql.Row) error {
	if err := row.Scan(&e.id, &e.name, &e.description, &e.isArchived, &e.createdAt, &e.updatedAt); err != nil {
		log.Println("error: ", err)
		return errorext.BuildDBError(err)
	}

	return nil
}

func (e *productEntity) scanRows(rows *sql.Rows) ([]productEntity, error) {
	d := []productEntity{}

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var p productEntity

		// fmt.Printf("Pointer: %p\n", &e)
		if err := rows.Scan(&p.id, &p.name, &p.description, &p.isArchived, &p.createdAt, &p.updatedAt); err != nil {
			log.Println("error: ", err)
			return nil, errorext.BuildDBError(err)
		}

		d = append(d, p)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Println("error: ", err)
		return nil, errorext.BuildDBError(err)
	}

	return d, nil
}
