package dto

type CreateProductRequest struct {
	TitleFa          string `json:"title_fa"`
	TitleEn          string `json:"title_en,omitempty"`
	Slug             string `json:"slug,omitempty"`
	Description      string `json:"description,omitempty"`
	BrandID          int64  `json:"brand_id"`
	CategoryID       int64  `json:"category_id"`
	OwnerType        string `json:"owner_type,omitempty"`
	OwnerID          *int64 `json:"owner_id,omitempty"`
	IsOriginal       *bool  `json:"is_original,omitempty"`
	MetaTitle        string `json:"meta_title,omitempty"`
	MetaDescription  string `json:"meta_description,omitempty"`
	IndexImageFileID *int64 `json:"index_image_file_id,omitempty"`
}

type UpdateProductRequest struct {
	TitleFa          *string `json:"title_fa,omitempty"`
	TitleEn          *string `json:"title_en,omitempty"`
	Slug             *string `json:"slug,omitempty"`
	Description      *string `json:"description,omitempty"`
	BrandID          *int64  `json:"brand_id,omitempty"`
	CategoryID       *int64  `json:"category_id,omitempty"`
	MetaTitle        *string `json:"meta_title,omitempty"`
	MetaDescription  *string `json:"meta_description,omitempty"`
	IndexImageFileID *int64  `json:"index_image_file_id,omitempty"`
}

type UpdateSEORequest struct {
	MetaTitle      *string `json:"meta_title,omitempty"`
	MetaDescription *string `json:"meta_description,omitempty"`
}
