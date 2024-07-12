package errorext

type CustomError struct {
	code int
	err  error
	// optional for addtional errors
	// might be helpful with logging
	additionalErrData map[string]any
}

func NewCustomError(code int, err error) *CustomError {
	return &CustomError{code: code, err: err}
}

func (e CustomError) Error() string {
	return e.err.Error()
}

func (e CustomError) Code() int {
	return e.code
}

func (e CustomError) Err() error {
	return e.err
}

func (e CustomError) AdditionalErrData() map[string]any {
	return e.additionalErrData
}

func (e *CustomError) SetAdditionalErrData(addErrData map[string]any) {
	e.additionalErrData = addErrData
}
