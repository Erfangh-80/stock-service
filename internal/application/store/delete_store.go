package store

import (
	"stock-service/internal/domain/store"
)

type DeleteStoreInput struct {
	ID int64
}

type DeleteStoreUseCase interface {
	Execute(input DeleteStoreInput) error
}

type deleteStoreInteractor struct {
	repo store.Repository
}

func NewDeleteStoreUseCase(repo store.Repository) DeleteStoreUseCase {
	return &deleteStoreInteractor{repo: repo}
}

func (uc *deleteStoreInteractor) Execute(input DeleteStoreInput) error {
	s, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return err
	}
	if s == nil {
		return store.ErrStoreNotFound
	}

	return uc.repo.Delete(input.ID)
}
