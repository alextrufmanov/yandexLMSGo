package calc

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyExpression   = errors.New("invalid input, expression is empty")
	ErrDivisionByZero    = errors.New("invalid input, division by zero")
	ErrInvalidExpression = errors.New(fmt.Sprintf("invalid input, expression"))
)

func ErrUnexpected(item string) error {
	return errors.New(fmt.Sprintf("invalid input, unexpected '%s'", item))
}
