package service

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrProductNotFound = errors.New("product not found")
)
