package app

import "fmt"

type ErrAppNotFound struct {
	Name string
}

//func NewErrAppNotFound(name string) error {
//	return ErrAppNotFound{Name: name}
//}

func (e ErrAppNotFound) Error() string {
	return fmt.Sprintf("app %s not found", e.Name)
}
