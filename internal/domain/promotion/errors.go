package promotion

import "errors"

var (
	ErrTitleRequired = errors.New("promotion title is required")
	ErrTitleTooLong  = errors.New("promotion title must not exceed 255 characters")
	ErrPromotionNotFound = errors.New("promotion not found")
)
