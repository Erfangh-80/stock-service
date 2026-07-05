package storewarehouselink

type RelationType string

const RelationTypePrimary RelationType = "primary"

// TODO: no created_at on schema — intentional for a junction table?
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

func (swl *StoreWarehouseLink) ChangeRelationType(rt RelationType) {
	swl.RelationType = rt
}
