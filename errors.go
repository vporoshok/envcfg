package envcfg

import "fmt"

const (
	// InvalidObjectType returned if object passed to processing is not pointer to struct
	InvalidObjectType constantError = "expected pointer to struct"
	// InvalidFieldType returned if some field has unsupported type
	InvalidFieldType constantError = "unsupported type"
)

type constantError string

func (ce constantError) Error() string {

	return string(ce)
}

func (ce constantError) New(msg string) error {

	return causedError{
		err: ce,
		msg: msg,
	}
}

type causedError struct {
	err error
	msg string
}

func (ce causedError) Error() string {

	return fmt.Sprintf("%s: %s", ce.msg, ce.err.Error())
}

func (ce causedError) Cause() error {

	return ce.err
}
