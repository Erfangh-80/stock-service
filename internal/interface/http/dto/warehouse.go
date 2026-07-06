package dto

type CreateWarehouseRequest struct {
	CreatedByUserID int64  `json:"created_by_user_id"`
	WarehouseName   string `json:"warehouse_name"`
}

type UpdateVisibilityRequest struct {
	IsPublic bool `json:"is_public"`
}

type UpdateWarehouseContactRequest struct {
	Phone            *string `json:"phone"`
	ContactPhone     *string `json:"contact_phone"`
	CollectionMethod string  `json:"collection_method"`
}
