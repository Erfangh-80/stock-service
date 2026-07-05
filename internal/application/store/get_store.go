package store

import (
	"stock-service/internal/domain/store"
)

type GetStoreInput struct {
	ID int64
}

type GetStoreUseCase interface {
	Execute(input GetStoreInput) (*store.Store, error)
}

type getStoreInteractor struct {
	repo store.Repository
}

func NewGetStoreUseCase(repo store.Repository) GetStoreUseCase {
	return &getStoreInteractor{repo: repo}
}

func (uc *getStoreInteractor) Execute(input GetStoreInput) (*store.Store, error) {
	s, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, store.ErrStoreNotFound
	}
	return s, nil
}
