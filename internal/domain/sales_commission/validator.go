package salescommission

func ValidateRatePercent(rate float64) error {
	if rate < 0 || rate > 100 {
		return ErrInvalidRatePercent
	}
	return nil
}

func ValidateMinPrice(price float64) error {
	if price <= 0 {
		return ErrInvalidMinPrice
	}
	return nil
}

func ValidateMaxPrice(min, max float64) error {
	if max <= min {
		return ErrInvalidMaxPrice
	}
	return nil
}
