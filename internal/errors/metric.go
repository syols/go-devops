package errors

import (
	"fmt"
)

type TypeNameError struct {
	info string
}

type HashSumError struct {
	info string
}

func (e *TypeNameError) Error() string {
	return fmt.Sprintf("wrong type name: %s", e.info)
}

func (e *HashSumError) Error() string {
	return fmt.Sprintf("wrong hash sum: %s", e.info)
}

func NewTypeNameError(info string) error {
	return &TypeNameError{
		info: info,
	}
}

func NewHashSumError(info string) error {
	return &HashSumError{
		info: info,
	}
}
