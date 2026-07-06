package dto

type CreateWarehouseLinkRequest struct {
	StoreID     int64 `json:"store_id"`
	WarehouseID int64 `json:"warehouse_id"`
}

type ChangeRelationRequest struct {
	RelationType string `json:"relation_type"`
}
