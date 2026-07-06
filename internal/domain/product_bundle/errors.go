package productbundle

import "errors"

var (
	ErrInvalidProductID       = errors.New("product id must be greater than zero")
	ErrInvalidRelatedProductID = errors.New("related product id must be greater than zero")
	ErrSelfReference          = errors.New("product cannot reference itself")
	ErrInvalidBundleType      = errors.New("invalid bundle type")
	ErrBundleNotFound         = errors.New("product bundle not found")
)
