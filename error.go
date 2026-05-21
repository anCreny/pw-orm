package pworm

import (
	"fmt"
)

type Error struct {
	message string
	err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.err.Error(), e.message)
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Wrap(err error) {
	e.err = err
}

func (e *Error) Is(err error) bool {
	if e.err != nil && err != nil {
		return err == e.err
	}

	return false
}

func NewError(err error, message string) *Error {
	return &Error{
		message: message,
		err:     err,
	}
}
