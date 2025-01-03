package postgres

import (
	"database/sql"
	"log"

	"github.com/tanveerprottoy/backend-structure-go/pkg/errorext"
)

// a package private entity clone of the domain entity
// this type can have db specific data types like here
// it's sqlext.NullString
type userEntity struct {
	id         string         `db:"id"`
	name       string         `db:"name"`
	address    sql.NullString `db:"address"`
	isArchived bool           `db:"is_archived"`
	createdAt  int64          `db:"created_at"`
	updatedAt  int64          `db:"updated_at"`
}

// this will be used to create db entity from domain entity
func newUserEntity(name string, address *string, createdAt, updatedAt int64) *userEntity {
	// address can be nil in db
	// check for nil
	e := userEntity{
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}

	if address != nil {
		e.address = sql.NullString{String: *address, Valid: true}
	}

	return &e
}

func (e *userEntity) scanRow(row *sql.Row) error {
	if err := row.Scan(&e.id, &e.name, &e.address, &e.isArchived, &e.createdAt, &e.updatedAt); err != nil {
		log.Println("error: ", err)
		return errorext.BuildDBError(err)
	}
	return nil
}

func (e *userEntity) scanRows(rows *sql.Rows) ([]userEntity, error) {
	d := []userEntity{}
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var p userEntity
		// fmt.Printf("Pointer: %p\n", &e)
		if err := rows.Scan(&p.id, &p.name, &p.address, &p.isArchived, &p.createdAt, &p.updatedAt); err != nil {
			log.Println("error: ", err)
			return nil, errorext.BuildDBError(err)
		}
		d = append(d, p)
	}
	return d, nil
}
