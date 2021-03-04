package utils

import (
	"errors"
	"fmt"
)

func PrintErrJoint(errinfo string, err error) error {
	errTotal := errinfo + "|" + err.Error()
	fmt.Println(errTotal)
	return errors.New(errTotal)
}

func ErrJoint(errinfo string, err error) error {
	errTotal := errinfo + " | " + err.Error()
	return errors.New(errTotal)
}
