package store

import (
	"stock-service/internal/domain/store"
)

type UpdateStoreNameInput struct {
	StoreID int64
	Name    string
}

type UpdateStoreNameUseCase interface {
	Execute(input UpdateStoreNameInput) (*store.Store, error)
}

type updateStoreNameInteractor struct {
	repo store.Repository
}

func NewUpdateStoreNameUseCase(repo store.Repository) UpdateStoreNameUseCase {
	return &updateStoreNameInteractor{repo: repo}
}

func (uc *updateStoreNameInteractor) Execute(input UpdateStoreNameInput) (*store.Store, error) {
	s, err := uc.repo.FindByID(input.StoreID)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, store.ErrStoreNotFound
	}

	if err := s.UpdateName(input.Name); err != nil {
		return nil, err
	}

	if err := uc.repo.Save(s); err != nil {
		return nil, err
	}

	return s, nil
}
