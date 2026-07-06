package dto

type CreateProductTypeRequest struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	SortOrder int    `json:"sort_order,omitempty"`
}
