package response

import (
	"encoding/json"
	"net/http"

	"github.com/tanveerprottoy/backend-structure-go/pkg/constant"
	"github.com/tanveerprottoy/backend-structure-go/pkg/typesext"
)

type Response[T any] struct {
	Data T `json:"data"`
}

func NewResponse[T any](payload T) *Response[T] {
	return &Response[T]{Data: payload}
}

type ReadManyResponse[T any] struct {
	Items []T `json:"items"`
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type Error struct {
	Message string `json:"message"`
}

func makeError(message string) Error {
	return Error{Message: message}
}

// helper function to build multiple errors
// used by MakeMultiErrorResponse
func buildMultipleErrors(errors []error) []Error {
	var errs []Error
	for _, err := range errors {
		errs = append(errs, makeError(err.Error()))
	}

	return errs
}

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

func NewErrorResponse(typ typesext.ErrorType, errors []error) *ErrorResponse {
	switch typ {
	case constant.ErrorSingle:
		return &ErrorResponse{Errors: []Error{makeError(errors[0].Error())}}
	case constant.ErrorMultiple:
		return &ErrorResponse{Errors: buildMultipleErrors(errors)}
	default:
		return &ErrorResponse{Errors: []Error{makeError(constant.InternalServerError)}}
	}
}

func RespondError(w http.ResponseWriter, code int, payload any) (int, error) {
	w.WriteHeader(code)

	res, errs := json.Marshal(payload)
	if errs != nil {
		// log failed to marshal
		return w.Write([]byte(constant.InternalServerError))
	}

	return w.Write(res)
}

func Respond(w http.ResponseWriter, code int, payload any) (int, error) {
	res, err := json.Marshal(payload)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, ErrorResponse{Errors: []Error{makeError(constant.InternalServerError)}})
		return -1, err
	}

	w.WriteHeader(code)

	return w.Write(res)
}
