package validatorext

import (
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
			msg = err.Field() + " " + err.Tag()
			errors = append(errors, errorext.ValidationError{Name: err.Field(), Message: msg})
		}
	}

	return errors
}
