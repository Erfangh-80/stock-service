package dto

type CreateStoreCategoryRequest struct {
	StoreID    int64 `json:"store_id"`
	CategoryID int64 `json:"category_id"`
}
