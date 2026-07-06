package storewarehouselink

import "errors"

var (
	ErrLinkNotFound       = errors.New("warehouse link not found")
	ErrInvalidRelationType = errors.New("invalid relation type")
)
