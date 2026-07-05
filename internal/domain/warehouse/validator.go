package warehouse

func ValidateWarehouseName(name string) error {
	if name == "" {
		return ErrWarehouseNameRequired
	}
	if len(name) > 255 {
		return ErrWarehouseNameTooLong
	}
	return nil
}
