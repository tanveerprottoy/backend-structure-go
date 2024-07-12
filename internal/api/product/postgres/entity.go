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
type productEntity struct {
	ID          string            `db:"id"`
	Name        string            `db:"name"`
	Description sqlext.NullString `db:"description"`
	IsArchived  bool              `db:"is_archived"`
	CreatedAt   int64             `db:"created_at"`
	UpdatedAt   int64             `db:"updated_at"`
}

func newProductEntity(name string, description *string, createdAt, updatedAt int64) *productEntity {
	e := productEntity{
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	if description != nil {
		e.Description = sqlext.MakeNullString(*description, true)
	}

	return &e
}

func (e *productEntity) scanRow(row *sql.Row) error {
	if err := row.Scan(&e.ID, &e.Name, &e.Description, &e.IsArchived, &e.CreatedAt, &e.UpdatedAt); err != nil {
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
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.IsArchived, &p.CreatedAt, &p.UpdatedAt); err != nil {
			log.Println("error: ", err)
			return nil, errorext.BuildDBError(err)
		}
		d = append(d, p)
	}
	return d, nil
}
