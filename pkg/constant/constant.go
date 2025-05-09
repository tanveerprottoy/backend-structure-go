package constant

import "github.com/tanveerprottoy/backend-structure-go/pkg/typesext"

const DBDriverName = "pgx"

// server
const ServerReadTimeout = 30
const ServerReadHeaderTimeout = 10
const ServerWriteTimeout = 30

// request
const RequestTimeout = 120

// cors configs
var AllowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
var AllowedHeaders = []string{"*"}
var AllowCredentials = false

// api url patterns
const ApiPattern string = "/api"
const V1 = "/v1"
const ProductsPattern = "/products"
const UsersPattern = "/users"

const InternalServerError = "internal server error"
const BadRequest = "bad request"
const NotFound = "not found"
const GenericFailMessage = "failed to perform the operation"
const InvalidQueryParam = "the query parameter supplied is invalid"
const MissingRequiredPathParam = "missing required path parameter id"
const InvalidRequestBody = "the request body is invalid"

const RequestTimeoutMsg string = "request timed out"

const ParamId = "id"
const ParamPage = "page"
const ParamLimit = "limit"
const ParamIsArchived = "isArchived"
const ParamSortBy = "sortBy"

const (
	ErrorSingle     typesext.ErrorType = "single"
	ErrorMultiple   typesext.ErrorType = "multiple"
	ErrorValidation typesext.ErrorType = "validation"
)

const EmptyUUID = "00000000-0000-0000-0000-000000000000"

const FakeUUID = "17e55148-8a8e-411c-bc72-028aecc8a20c"
