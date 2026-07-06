package storeinterface

import (
	appstore "stock-service/internal/application/store"
	"stock-service/internal/domain/store"
	iface "stock-service/internal/interface"
)

type StoreResponse struct {
	ID                     int64          `json:"id"`
	UserID                 int64          `json:"user_id"`
	StoreName              string         `json:"store_name"`
	Status                 string         `json:"status"`
	AddressID              *int64         `json:"address_id,omitempty"`
	ContactPhone           *string        `json:"contact_phone,omitempty"`
	MediaAssets            map[string]any `json:"media_assets,omitempty"`
	IsCommissionApplicable bool           `json:"is_commission_applicable"`
	IsBulkSaleEnabled      bool           `json:"is_bulk_sale_enabled"`
	CreatedAt              string         `json:"created_at"`
}

type ListStoresResponse struct {
	Stores []StoreResponse `json:"stores"`
	Total  int             `json:"total"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
}

type CreateStoreParams struct {
	UserID    int64
	StoreName string
}

type UpdateContactParams struct {
	StoreID      int64
	ContactPhone *string
}

type UpdateNameParams struct {
	StoreID int64
	Name    string
}

type UpdateProfileParams struct {
	StoreID     int64
	AddressID   *int64
	MediaAssets map[string]any
}

type Adapter struct {
	create        appstore.CreateStoreUseCase
	get           appstore.GetStoreUseCase
	list          appstore.ListStoresUseCase
	toggleBulk    appstore.ToggleBulkSaleUseCase
	toggleComm    appstore.ToggleCommissionUseCase
	updateContact appstore.UpdateContactUseCase
	updateName    appstore.UpdateStoreNameUseCase
	updateProfile appstore.UpdateStoreProfileUseCase
	delete        appstore.DeleteStoreUseCase
}

func NewAdapter(
	create appstore.CreateStoreUseCase,
	get appstore.GetStoreUseCase,
	list appstore.ListStoresUseCase,
	toggleBulk appstore.ToggleBulkSaleUseCase,
	toggleComm appstore.ToggleCommissionUseCase,
	updateContact appstore.UpdateContactUseCase,
	updateName appstore.UpdateStoreNameUseCase,
	updateProfile appstore.UpdateStoreProfileUseCase,
	delete appstore.DeleteStoreUseCase,
) *Adapter {
	return &Adapter{
		create: create, get: get, list: list,
		toggleBulk: toggleBulk, toggleComm: toggleComm,
		updateContact: updateContact, updateName: updateName,
		updateProfile: updateProfile, delete: delete,
	}
}

func (a *Adapter) Create(params CreateStoreParams) (*StoreResponse, error) {
	result, err := a.create.Execute(appstore.CreateStoreInput{
		UserID: params.UserID, StoreName: params.StoreName,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Get(id int64) (*StoreResponse, error) {
	result, err := a.get.Execute(appstore.GetStoreInput{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

type ListStoresFilter struct {
	UserID *int64
	Status *string
	Page   int
	Limit  int
}

func (a *Adapter) List(filter ListStoresFilter) (*ListStoresResponse, error) {
	result, err := a.list.Execute(appstore.ListStoresInput{
		UserID: filter.UserID,
		Status: filter.Status,
		Page:   filter.Page,
		Limit:  filter.Limit,
	})
	if err != nil {
		return nil, mapError(err)
	}
	stores := make([]StoreResponse, len(result.Stores))
	for i, s := range result.Stores {
		stores[i] = *toResponse(s)
	}
	return &ListStoresResponse{
		Stores: stores,
		Total:  result.Total,
		Page:   result.Page,
		Limit:  result.Limit,
	}, nil
}

func (a *Adapter) ToggleBulkSale(storeID int64) (*StoreResponse, error) {
	result, err := a.toggleBulk.Execute(appstore.ToggleBulkSaleInput{StoreID: storeID})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) ToggleCommission(storeID int64) (*StoreResponse, error) {
	result, err := a.toggleComm.Execute(appstore.ToggleCommissionInput{StoreID: storeID})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) UpdateContact(params UpdateContactParams) (*StoreResponse, error) {
	result, err := a.updateContact.Execute(appstore.UpdateContactInput{
		StoreID: params.StoreID, ContactPhone: params.ContactPhone,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) UpdateName(params UpdateNameParams) (*StoreResponse, error) {
	result, err := a.updateName.Execute(appstore.UpdateStoreNameInput{
		StoreID: params.StoreID, Name: params.Name,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) UpdateProfile(params UpdateProfileParams) (*StoreResponse, error) {
	result, err := a.updateProfile.Execute(appstore.UpdateStoreProfileInput{
		StoreID:     params.StoreID,
		AddressID:   params.AddressID,
		MediaAssets: params.MediaAssets,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return toResponse(result), nil
}

func (a *Adapter) Delete(id int64) error {
	err := a.delete.Execute(appstore.DeleteStoreInput{ID: id})
	if err != nil {
		return mapError(err)
	}
	return nil
}

func mapError(err error) error {
	switch err {
	case store.ErrStoreNotFound:
		return iface.ErrNotFound
	case store.ErrStoreNameRequired, store.ErrStoreNameTooLong:
		return iface.ErrInvalidInput
	default:
		return iface.ErrInternal
	}
}

func toResponse(s *store.Store) *StoreResponse {
	return &StoreResponse{
		ID: s.ID, UserID: s.UserID, StoreName: s.StoreName,
		Status: string(s.Status), AddressID: s.AddressID,
		ContactPhone: s.ContactPhone, MediaAssets: s.MediaAssets,
		IsCommissionApplicable: s.IsCommissionApplicable,
		IsBulkSaleEnabled:      s.IsBulkSaleEnabled,
		CreatedAt:              s.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
