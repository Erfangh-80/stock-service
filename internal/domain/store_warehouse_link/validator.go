package storewarehouselink

var validRelationTypes = map[RelationType]bool{
	RelationTypePrimary: true,
}

func ValidateRelationType(rt RelationType) error {
	if !validRelationTypes[rt] {
		return ErrInvalidRelationType
	}
	return nil
}
