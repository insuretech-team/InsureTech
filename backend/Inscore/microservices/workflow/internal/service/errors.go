package service

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInvalidTransition = errors.New("invalid state transition")
	ErrFeatureUnavailable = errors.New("feature unavailable")
)
