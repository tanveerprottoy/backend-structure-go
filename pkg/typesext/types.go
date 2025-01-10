// package typesext contains custom/derived/union types
package typesext

type Numeric interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type NumStrBool interface {
	Numeric | string | bool
}

// ContextKey represents the key of a context
type ContextKey string

// ErrorType represents the type of the error
type ErrorType string
