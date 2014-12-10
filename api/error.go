package api

import "fmt"

type Error struct {
	Status  int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("API error (%d): %s", e.Status, e.Message)
}

func newError(status int, body []byte) *Error {
	return &Error{Status: status, Message: string(body)}
}
