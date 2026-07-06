package brand

import "errors"

var (
	ErrNameRequired = errors.New("brand name is required")
	ErrSlugRequired = errors.New("brand slug is required")
	ErrBrandNotFound = errors.New("brand not found")
)
