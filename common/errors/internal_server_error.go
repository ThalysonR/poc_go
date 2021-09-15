package errors

type InternalServerError struct {
}

func NewInternalServerError() *InternalServerError {
	return &InternalServerError{}
}

func (i *InternalServerError) Error() string {
	return "Internal Server Error"
}
