package errors

type ErrInsufficientData struct {
	Message	string
}

func (e *ErrInsufficientData) Error() string {
	return e.Message
}
