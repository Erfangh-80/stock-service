package referenceprice

import "errors"

var (
	ErrInvalidReferencePrice = errors.New("reference price must be greater than zero")
	ErrReferencePriceNotFound = errors.New("reference price not found")
)
