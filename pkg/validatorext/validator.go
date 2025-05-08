package validatorext

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/tanveerprottoy/backend-structure-go/pkg/errorext"
)

type Validater interface {
	Validate(v any) []error
}

type Validator struct {
	validate *validator.Validate
}

func NewValidator(validate *validator.Validate) *Validator {
	v := &Validator{validate: validate}
	v.registerTagNameFunc()
	return v
}

// RegisterTagNameFunc configures validator to use
// defined json name to use as struct field name
func (v *Validator) registerTagNameFunc() {
	v.validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		n := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]

		// skip if tag key says it should be ignored
		if n == "-" {
			return ""
		}

		return n
	})
}

// The caller must pass the address for the param val
// the param val must be a struct
func (v *Validator) Validate(val any) []error {
	var errors []error

	// validate request body
	err := v.validate.Struct(val)
	if err != nil {
		// Struct is invalid
		var msg string

		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			// need to handle some tags
			case "oneof":
				// set oneof specific error message
				msg = fmt.Sprintf("field %s has incorrect value expected one of: %s", err.Field(), err.Param())
			case "gt":
				// set gt specific error message
				msg = fmt.Sprintf("field %s must be greater than %s", err.Field(), err.Param())
			case "gte":
				// set gte specific error message
				msg = fmt.Sprintf("field %s must be greater or equal %s", err.Field(), err.Param())
			case "min":
				// set min specific error message
				msg = fmt.Sprintf("field %s must be minimum %s", err.Field(), err.Param())
			case "len":
				// set len specific error message
				msg = fmt.Sprintf("field %s length must be exactly %s", err.Field(), err.Param())
			default:
				msg = err.Field() + " " + err.Tag()
			}

			errors = append(errors, errorext.ValidationError{Name: err.Field(), Message: msg})
		}
	}

	return errors
}
