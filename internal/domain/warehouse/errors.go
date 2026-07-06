package warehouse

import "errors"

var (
	ErrWarehouseNameRequired         = errors.New("warehouse name is required")
	ErrWarehouseNameTooLong          = errors.New("warehouse name must not exceed 255 characters")
	ErrWarehouseNotFound             = errors.New("warehouse not found")
	ErrInvalidCollectionMethod       = errors.New("invalid collection method")
	ErrWarehouseAddressIDNotPositive = errors.New("address id must be positive")
)
