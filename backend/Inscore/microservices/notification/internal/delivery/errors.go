package delivery

import "errors"

type permanentError struct {
	err error
}

func (e *permanentError) Error() string {
	return e.err.Error()
}

func (e *permanentError) Unwrap() error {
	return e.err
}

func Permanent(err error) error {
	if err == nil {
		return nil
	}
	var target *permanentError
	if errors.As(err, &target) {
		return err
	}
	return &permanentError{err: err}
}

func IsPermanent(err error) bool {
	var target *permanentError
	return errors.As(err, &target)
}
