package referenceprice

func ValidatePrice(price float64) error {
	if price <= 0 {
		return ErrInvalidReferencePrice
	}
	return nil
}
