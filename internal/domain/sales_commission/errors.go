package salescommission

import "errors"

var (
	ErrInvalidRatePercent = errors.New("commission rate percent must be between 0 and 100")
	ErrInvalidMinPrice    = errors.New("minimum price must be greater than zero")
	ErrInvalidMaxPrice    = errors.New("maximum price must be greater than minimum price")
	ErrInvalidMinQty      = errors.New("minimum quantity must be greater than or equal to zero")
	ErrCommissionNotFound = errors.New("sales commission not found")
	ErrRuleNotFound       = errors.New("category commission rule not found")
)
