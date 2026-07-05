package store

import (
	"stock-service/internal/domain/store"
)

type UpdateStoreProfileInput struct {
	StoreID     int64
	AddressID   *int64
	MediaAssets map[string]any
}

type UpdateStoreProfileUseCase interface {
	Execute(input UpdateStoreProfileInput) (*store.Store, error)
}

type updateStoreProfileInteractor struct {
	repo store.Repository
}

func NewUpdateStoreProfileUseCase(repo store.Repository) UpdateStoreProfileUseCase {
	return &updateStoreProfileInteractor{repo: repo}
}

func (uc *updateStoreProfileInteractor) Execute(input UpdateStoreProfileInput) (*store.Store, error) {
	s, err := uc.repo.FindByID(input.StoreID)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, store.ErrStoreNotFound
	}

	s.AddressID = input.AddressID
	if input.MediaAssets != nil {
		s.MediaAssets = input.MediaAssets
	}

	if err := uc.repo.Save(s); err != nil {
		return nil, err
	}

	return s, nil
}
