package storewarehouselink

type RelationType string

const RelationTypePrimary RelationType = "primary"

type StoreWarehouseLink struct {
	ID           int64
	StoreID      int64
	WarehouseID  int64
	RelationType RelationType
}

func NewStoreWarehouseLink(storeID, warehouseID int64) *StoreWarehouseLink {
	return &StoreWarehouseLink{
		StoreID:      storeID,
		WarehouseID:  warehouseID,
		RelationType: RelationTypePrimary,
	}
}

func (swl *StoreWarehouseLink) ChangeRelationType(rt RelationType) error {
	if err := ValidateRelationType(rt); err != nil {
		return err
	}
	swl.RelationType = rt
	return nil
}
