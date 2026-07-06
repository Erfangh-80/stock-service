package promotion

import "time"

func ValidateTitle(title string) error {
	if title == "" {
		return ErrTitleRequired
	}
	if len(title) > 255 {
		return ErrTitleTooLong
	}
	return nil
}

func ValidateDiscountType(dt DiscountType) error {
	if dt == "" {
		return ErrDiscountTypeRequired
	}
	if dt != DiscountTypePercentage && dt != DiscountTypeFixedAmount {
		return ErrInvalidDiscountType
	}
	return nil
}

func ValidateDiscountValue(dt DiscountType, value float64) error {
	if value <= 0 {
		return ErrDiscountValueRequired
	}
	if dt == DiscountTypePercentage && value > 100 {
		return ErrDiscountValueTooHigh
	}
	return nil
}

func ValidateCouponCode(code string) error {
	if len(code) > 50 {
		return ErrInvalidCouponCode
	}
	return nil
}

func ValidatePromotionDates(start, end *time.Time) error {
	if start != nil && end != nil {
		if !end.After(*start) {
			return ErrInvalidPromotionDates
		}
	}
	return nil
}
