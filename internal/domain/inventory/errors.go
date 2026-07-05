package inventory

import "errors"

var (
	ErrBasePriceRequired        = errors.New("base price is required")
	ErrInvalidBasePrice         = errors.New("base price must be greater than zero")
	ErrInvalidFinalPrice        = errors.New("final price must be greater than zero")
	ErrInvalidQuantity          = errors.New("quantity must be greater than or equal to zero")
	ErrInvalidMinOrderQty       = errors.New("minimum order quantity must be greater than zero")
	ErrInvalidMaxOrderQty       = errors.New("maximum order quantity must be greater than or equal to minimum order quantity")
	ErrInvalidPromotionDates    = errors.New("promotion end date must be after start date")
	ErrPromotionAlreadyApplied  = errors.New("promotion already applied to this inventory item")
	ErrNoActivePromotion        = errors.New("no active promotion to remove")
	ErrInventoryNotFound        = errors.New("inventory not found")
	ErrVendorSaleStatusTransition = errors.New("invalid vendor sale status transition")
	ErrSystemSaleStatusTransition = errors.New("invalid system sale status transition")
	ErrInvalidScheduledDate     = errors.New("scheduled delivery date must be in the future")
	ErrInsufficientStock        = errors.New("insufficient stock to reserve quantity")
)
