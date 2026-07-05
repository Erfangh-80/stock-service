package warehouse

import "errors"

var (
	ErrWarehouseNameRequired = errors.New("warehouse name is required")
	ErrWarehouseNameTooLong  = errors.New("warehouse name must not exceed 255 characters")
)
