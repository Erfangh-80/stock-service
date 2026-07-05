package store

import "errors"

var (
	ErrStoreNameRequired = errors.New("store name is required")
	ErrStoreNameTooLong  = errors.New("store name must not exceed 255 characters")
	ErrStoreNotFound     = errors.New("store not found")
)
