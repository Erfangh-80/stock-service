package storeallowedcategory

func ValidateSupportNote(note string) error {
	if len(note) > 500 {
		return ErrSupportNoteTooLong
	}
	return nil
}
