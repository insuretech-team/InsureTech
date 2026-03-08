package service

import "errors"

var (
    ErrInvalidArgument   = errors.New("invalid argument")
    ErrNotFound          = errors.New("payment not found")
    ErrInvalidTransition = errors.New("invalid payment state transition")
    ErrPaymentFailed     = errors.New("payment failed")
    ErrNotImplemented    = errors.New("payment feature not implemented")
)
