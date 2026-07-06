package promotion

import "errors"

var (
	ErrTitleRequired              = errors.New("promotion title is required")
	ErrTitleTooLong               = errors.New("promotion title must not exceed 255 characters")
	ErrPromotionNotFound          = errors.New("promotion not found")
	ErrDiscountTypeRequired       = errors.New("promotion discount type is required")
	ErrInvalidDiscountType        = errors.New("invalid promotion discount type")
	ErrDiscountValueRequired      = errors.New("promotion discount value is required")
	ErrInvalidDiscountValue       = errors.New("promotion discount value must be greater than zero")
	ErrDiscountValueTooHigh       = errors.New("percentage discount value must not exceed 100")
	ErrInvalidPromotionDates      = errors.New("promotion end date must be after start date")
	ErrPromotionExpired           = errors.New("promotion has expired")
	ErrPromotionNotActive         = errors.New("promotion is not active")
	ErrPromotionNotStarted        = errors.New("promotion has not started yet")
	ErrPromotionUsageLimitExceeded = errors.New("promotion usage limit exceeded")
	ErrPromotionBudgetExceeded    = errors.New("promotion budget exceeded")
	ErrInvalidCouponCode          = errors.New("coupon code must not exceed 50 characters")
	ErrCouponCodeAlreadyExists    = errors.New("coupon code already exists")
	ErrIneligibleStore            = errors.New("promotion is not eligible for this store")
	ErrIneligibleCategory         = errors.New("promotion is not eligible for this category")
	ErrIneligibleProduct          = errors.New("promotion is not eligible for this product")
	ErrIneligibleUser             = errors.New("promotion is not eligible for this user")
	ErrInvalidDateFormat          = errors.New("invalid date format, use RFC3339")
)
