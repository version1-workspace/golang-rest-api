package model

import "fmt"

func IsErrorNotFound(err error) bool {
	_, ok := err.(*ErrNotFound)
	return ok
}

func NewErrorNotFound(name string, id int) error {
	return &ErrNotFound{resourceName: name, id: id}
}

type ErrNotFound struct {
	resourceName string
	id           int
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("not found: %s %d", e.resourceName, e.id)
}
