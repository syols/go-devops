package errors

import "fmt"

type ValueNotFound struct {
	info string
}

func (e *ValueNotFound) Error() string {
	return fmt.Sprintf("value not found: %s", e.info)
}

func NewValueNotFound(info string) error {
	return &TypeNameError{
		info: info,
	}
}
