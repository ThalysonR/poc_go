package errors

import "fmt"

type ErrNotFound struct {
	msg string
}

func NewErrNotFound(msg string, args ...interface{}) *ErrNotFound {
	return &ErrNotFound{
		msg: fmt.Sprintf(msg, args...),
	}
}

func (e *ErrNotFound) Error() string {
	return e.msg
}
