package inventory

import "time"

func ValidateBasePrice(price float64) error {
	if price <= 0 {
		return ErrInvalidBasePrice
	}
	return nil
}

func ValidateFinalPrice(price float64) error {
	if price <= 0 {
		return ErrInvalidFinalPrice
	}
	return nil
}

func ValidateInstantQty(qty int) error {
	if qty < 0 {
		return ErrInvalidQuantity
	}
	return nil
}

func ValidateMinOrderQty(qty int) error {
	if qty < 1 {
		return ErrInvalidMinOrderQty
	}
	return nil
}

func ValidateMaxOrderQty(min, max int) error {
	if max < min {
		return ErrInvalidMaxOrderQty
	}
	return nil
}

func ValidatePromotionDates(start, end time.Time) error {
	if !end.After(start) {
		return ErrInvalidPromotionDates
	}
	return nil
}

func ValidateScheduledDeliveryDate(date string) error {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ErrInvalidScheduledDate
	}
	if !t.After(time.Now().Add(-24 * time.Hour)) {
		return ErrInvalidScheduledDate
	}
	return nil
}
