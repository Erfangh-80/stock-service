package dto

type CreateProductBundleRequest struct {
	RelatedProductID int32  `json:"related_product_id"`
	Type             string `json:"type"`
	SortOrder        int    `json:"sort_order,omitempty"`
}
