package errorext

type ValidationError struct {
	Name    string
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}
