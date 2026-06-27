package service

import "errors"

var (
	ErrOrderNotFound   = errors.New("order not found")
	ErrProductNotFound = errors.New("product not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidInput    = errors.New("invalid input")
)
