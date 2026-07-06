package warehouse

var validCollectionMethods = map[string]bool{
	"pickup":    true,
	"delivery":  true,
	"both":      true,
}

func ValidateWarehouseName(name string) error {
	if name == "" {
		return ErrWarehouseNameRequired
	}
	if len(name) > 255 {
		return ErrWarehouseNameTooLong
	}
	return nil
}

func ValidateCollectionMethod(method string) error {
	if !validCollectionMethods[method] {
		return ErrInvalidCollectionMethod
	}
	return nil
}

func ValidateAddressID(id int64) error {
	if id <= 0 {
		return ErrWarehouseAddressIDNotPositive
	}
	return nil
}
