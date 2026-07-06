package producttype

import "errors"

var (
	ErrInvalidProductID = errors.New("product id must be greater than zero")
	ErrNameRequired     = errors.New("type name is required")
	ErrValueRequired    = errors.New("type value is required")
	ErrTypeNotFound     = errors.New("product type not found")
)
