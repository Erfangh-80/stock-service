package productimage

import "errors"

var (
	ErrInvalidProductID = errors.New("product id must be greater than zero")
	ErrInvalidFileID    = errors.New("file id must be greater than zero")
	ErrImageNotFound    = errors.New("product image not found")
)
