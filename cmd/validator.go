package cmd

import (
	"errors"
)

var cmdValidator cmdOptionValidator

type cmdOptionValidator struct {
	cmdErrors []error
}

func (v *cmdOptionValidator) valid() bool {
	return len(v.cmdErrors) == 0
}

func (v *cmdOptionValidator) check(condition bool, failureMsg string) {
	err := errors.New(failureMsg)
	if !condition {
		v.cmdErrors = append(v.cmdErrors, err)
	}
}
