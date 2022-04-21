package cluster

import "fmt"

type NotFound struct {
	Name string
	Err  error
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("cluster '%s' not found", e.Name)
}

func (e *NotFound) Unwrap() error { return e.Err }

type AlreadyRegistered struct {
	Name string
	Err  error
}

func (e *AlreadyRegistered) Error() string {
	return fmt.Sprintf("cluster '%s' already registered", e.Name)
}

func (e *AlreadyRegistered) Unwrap() error { return e.Err }
