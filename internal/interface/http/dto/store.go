package dto

type CreateStoreRequest struct {
	UserID    int64  `json:"user_id"`
	StoreName string `json:"store_name"`
}

type UpdateStoreContactRequest struct {
	ContactPhone *string `json:"contact_phone"`
}

type UpdateStoreNameRequest struct {
	Name string `json:"name"`
}

type UpdateStoreProfileRequest struct {
	AddressID   *int64         `json:"address_id,omitempty"`
	MediaAssets map[string]any `json:"media_assets,omitempty"`
}
