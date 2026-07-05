package product

import "errors"

var (
	ErrInvalidProductID = errors.New("product id must be greater than zero")
	ErrProductNotFound  = errors.New("product not found")
	ErrTitleFaRequired  = errors.New("persian title is required")
	ErrInvalidBrandID   = errors.New("brand id must be greater than zero")
	ErrInvalidCategoryID = errors.New("category id must be greater than zero")
)
