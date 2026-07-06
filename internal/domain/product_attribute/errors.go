package productattribute

import "errors"

var (
	ErrInvalidProductID = errors.New("product id must be greater than zero")
	ErrKeyRequired      = errors.New("attribute key is required")
	ErrAttributeNotFound = errors.New("product attribute not found")
)
