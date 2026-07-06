package category

func ValidateCategoryName(name string) error {
	if name == "" {
		return ErrNameRequired
	}
	return nil
}

func ValidateCategorySlug(slug string) error {
	if slug == "" {
		return ErrSlugRequired
	}
	return nil
}
