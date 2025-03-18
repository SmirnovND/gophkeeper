package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrInsufficientFunds = errors.New("insufficient funds")

type Error struct {
	Message   string
	CodeValue int
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Code() int {
	return e.CodeValue
}
