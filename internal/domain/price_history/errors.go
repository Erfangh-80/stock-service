package pricehistory

import "errors"

var (
	ErrInvalidProductID  = errors.New("product id must be greater than zero")
	ErrInvalidPrice      = errors.New("price must not be negative")
	ErrChangedByRequired = errors.New("changed by is required")
)
