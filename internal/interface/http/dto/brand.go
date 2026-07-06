package dto

type CreateBrandRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type UpdateBrandRequest struct {
	Name *string `json:"name,omitempty"`
	Slug *string `json:"slug,omitempty"`
}
