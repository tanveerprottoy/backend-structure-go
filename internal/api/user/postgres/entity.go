package postgres

import (
	"database/sql"
	"log"

	"github.com/tanveerprottoy/backend-structure-go/pkg/errorext"
	"github.com/tanveerprottoy/backend-structure-go/pkg/sqlext"
)

// a package private entity clone of the domain entity
// this type can have db specific data types like here
// it's sqlext.NullString
type userEntity struct {
	ID         string            `db:"id"`
	Name       string            `db:"name"`
	Address    sqlext.NullString `db:"address"`
	IsArchived bool              `db:"is_archived"`
	CreatedAt  int64             `db:"created_at"`
	UpdatedAt  int64             `db:"updated_at"`
}

func newUserEntity(name string, address *string, createdAt, updatedAt int64) *userEntity {
	e := userEntity{
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	if address != nil {
		e.Address = sqlext.MakeNullString(*address, true)
	}

	return &e
}

func (e *userEntity) scanRow(row *sql.Row) error {
	if err := row.Scan(&e.ID, &e.Name, &e.Address, &e.IsArchived, &e.CreatedAt, &e.UpdatedAt); err != nil {
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
		if err := rows.Scan(&p.ID, &p.Name, &p.Address, &p.IsArchived, &p.CreatedAt, &p.UpdatedAt); err != nil {
			log.Println("error: ", err)
			return nil, errorext.BuildDBError(err)
		}
		d = append(d, p)
	}
	return d, nil
}
