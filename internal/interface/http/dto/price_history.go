package dto

type CreatePriceHistoryRequest struct {
	OldPrice    float64 `json:"old_price"`
	NewPrice    float64 `json:"new_price"`
	ChangedBy   string  `json:"changed_by"`
	Description *string `json:"description,omitempty"`
}
