package storeallowedcategory

import "errors"

var (
	ErrStoreCategoryNotFound = errors.New("store category not found")
	ErrCategoryNotFound      = errors.New("category not found")
	ErrSupportNoteTooLong    = errors.New("support note is too long")
)
