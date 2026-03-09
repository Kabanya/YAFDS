package errors

import "errors"

var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrCourierNotFound  = errors.New("courier not found")
	ErrOrderNotFound    = errors.New("order not found")
)
