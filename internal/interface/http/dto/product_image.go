package dto

type CreateProductImageRequest struct {
	FileID    int64 `json:"file_id"`
	SortOrder int   `json:"sort_order,omitempty"`
}
