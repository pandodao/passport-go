package eip4361

import "fmt"

type Err struct {
	err error
	Key Key
}

func (e *Err) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Key.String(), e.err.Error())
}

func (e *Err) Unwrap() error {
	return e.err
}


