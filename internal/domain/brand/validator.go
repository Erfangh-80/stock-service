package brand

func ValidateBrandName(name string) error {
	if name == "" {
		return ErrNameRequired
	}
	return nil
}

func ValidateBrandSlug(slug string) error {
	if slug == "" {
		return ErrSlugRequired
	}
	return nil
}
