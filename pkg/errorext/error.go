package errorext

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tanveerprottoy/backend-structure-go/pkg/constant"
)

const (
	// case_not_found
	SQLCodeNotFound = "20000"
	// no_data
	SQLCodeNoData = "02000"
	// invalid_text_representation
	SQLCodeInvalidTextRepresentation = "22P02"
	// undefined_function
	SQLCodeUndefinedFunction = "42883"
	// undefined_table
	SQLCodeUndefinedTable = "42P01"
	// undefined_parameter
	SQLCodeUndefinedParam = "42P02"
	// invalid_column_reference
	SQLInvalidColumnReference = "42P10"
	// unique violation
	SQLCodeUniqueViolation = "23505"
)

var ErrNotFound = errors.New(constant.NotFound)
var ErrInternalServer = errors.New(constant.InternalServerError)

func BuildCustomError(err error) error {
	var customErr CustomError
	ok := errors.As(err, &customErr)
	if ok {
		return &customErr
	}
	
	// return custom error with code
	return NewCustomError(http.StatusInternalServerError, err)
}

func ParseCustomError(err error) *CustomError {
	var customErr CustomError
	ok := errors.As(err, &customErr)
	if ok {
		return &customErr
	}

	// return custom error with code
	return NewCustomError(http.StatusInternalServerError, err)
}

func ParseJSONError(err error) error {
	if err, ok := err.(*json.UnmarshalTypeError); ok {
		return errors.New("invalid type for " + err.Field)
	}

	return errors.New(constant.InvalidRequestBody)
}

func BuildDBError(err error) error {
	// check if it's an sql error
	switch err {
	case sql.ErrNoRows:
		return NewCustomError(http.StatusNotFound, ErrNotFound)
	case sql.ErrTxDone:
		cerr := NewCustomError(http.StatusInternalServerError, ErrInternalServer)
		cerr.SetAdditionalErrData(map[string]any{"code": "", "message": "transaction already closed", "detail": ""})
		return cerr
	}
	return NewCustomError(http.StatusInternalServerError, ErrInternalServer)
}
