package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"

	"github.com/tanveerprottoy/backend-structure-go/pkg/errorext"
)

// a package private entity clone of the domain entity
// this type can have db specific data types like here
// it's sqlext.NullString
type userEntity struct {
	// this struct fields must be exported
	// so that the reflection can access them
	// which is done in the method scanMany
	Id         string         `db:"id"`
	Name       string         `db:"name"`
	Address    sql.NullString `db:"address"`
	IsArchived bool           `db:"is_archived"`
	CreatedAt  int64          `db:"created_at"`
	UpdatedAt  int64          `db:"updated_at"`
}

// this will be used to create db entity from domain entity
func newUserEntity(name string, address *string, createdAt, updatedAt int64) *userEntity {
	// address can be nil in db
	// check for nil
	e := userEntity{
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	if address != nil {
		e.Address = sql.NullString{String: *address, Valid: true}
	}

	return &e
}

func (e *userEntity) scanRow(row *sql.Row) error {
	if err := row.Scan(&e.Id, &e.Name, &e.Address, &e.IsArchived, &e.CreatedAt, &e.UpdatedAt); err != nil {
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
		if err := rows.Scan(&p.Id, &p.Name, &p.Address, &p.IsArchived, &p.CreatedAt, &p.UpdatedAt); err != nil {
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

func (e *userEntity) scanMany(rows *sql.Rows) ([]userEntity, error) {
	d := make([]userEntity, 0)

	columnNames, err := rows.Columns()
	if err != nil {
		// handle err
		return d, nil
	}

	// put the column names in a map
	columnIndexes := make(map[string]int, len(columnNames))
	for i, v := range columnNames {
		columnIndexes[v] = i
	}

	for rows.Next() {
		entity := userEntity{}
		pointers := make([]any, len(columnNames))
		// pointers array's index
		j := 0

		val := reflect.ValueOf(&entity).Elem()

		for i := 0; i < val.NumField(); i++ {
			f := val.Field(i)
			ft := val.Type().Field(i)

			// Check if the field is valid
			if !f.IsValid() {
				return d, fmt.Errorf("field %s is not valid", ft.Name)
			}

			// Check if the field is addressable
			if !f.CanAddr() {
				return d, fmt.Errorf("field %s is not addressable", ft.Name)
			}

			// get the tag of the field
			tag := ft.Tag.Get("db")
			if _, ok := columnIndexes[tag]; ok {
				// add the pointer to the slice
				pointers[j] = f.Addr().Interface()

				// increment the pointers arr index
				j++
			}
			// Check if the field can be set
			/* if !fv.CanSet() {
				return []User{}, fmt.Errorf("field %s is not settable", f.Name)
			} */
		}

		err := rows.Scan(pointers...)
		if err != nil {
			// handle err
			return d, fmt.Errorf("error: %w", err)
		}

		d = append(d, entity)
	}

	if err := rows.Err(); err != nil {
		log.Println("error: ", err)
		return nil, errorext.BuildDBError(err)
	}

	return d, nil
}
