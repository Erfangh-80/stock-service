package category

import "errors"

var (
	ErrNameRequired  = errors.New("category name is required")
	ErrSlugRequired  = errors.New("category slug is required")
	ErrCategoryNotFound = errors.New("category not found")
)
