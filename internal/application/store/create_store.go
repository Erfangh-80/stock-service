package store

import (
	"stock-service/internal/domain/store"
)

type CreateStoreInput struct {
	UserID    int64
	StoreName string
}

type CreateStoreUseCase interface {
	Execute(input CreateStoreInput) (*store.Store, error)
}

type createStoreInteractor struct {
	repo store.Repository
}

func NewCreateStoreUseCase(repo store.Repository) CreateStoreUseCase {
	return &createStoreInteractor{repo: repo}
}

func (uc *createStoreInteractor) Execute(input CreateStoreInput) (*store.Store, error) {
	s, err := store.NewStore(input.UserID, input.StoreName)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(s); err != nil {
		return nil, err
	}

	return s, nil
}
