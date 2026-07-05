package store

func ValidateStoreName(name string) error {
	if name == "" {
		return ErrStoreNameRequired
	}
	if len(name) > 255 {
		return ErrStoreNameTooLong
	}
	return nil
}
