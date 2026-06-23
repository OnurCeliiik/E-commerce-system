package service

import "errors"

var (
	ErrInventoryNotFound = errors.New("inventory not found")
	ErrInvalidInput      = errors.New("invalid input")
)
