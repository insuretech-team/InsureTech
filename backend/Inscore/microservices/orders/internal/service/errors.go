package service

import "errors"

var (
	ErrInvalidArgument   = errors.New("invalid argument")
	ErrNotFound          = errors.New("not found")
	ErrAlreadyExists     = errors.New("already exists")
	ErrInvalidTransition = errors.New("invalid state transition")
	ErrPaymentFailed     = errors.New("payment failed")
	ErrNotImplemented    = errors.New("not implemented")
)
