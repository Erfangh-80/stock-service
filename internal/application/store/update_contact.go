package store

import (
	"stock-service/internal/domain/store"
)

type UpdateContactInput struct {
	StoreID      int64
	ContactPhone *string
}

type UpdateContactUseCase interface {
	Execute(input UpdateContactInput) (*store.Store, error)
}

type updateContactInteractor struct {
	repo store.Repository
}

func NewUpdateContactUseCase(repo store.Repository) UpdateContactUseCase {
	return &updateContactInteractor{repo: repo}
}

func (uc *updateContactInteractor) Execute(input UpdateContactInput) (*store.Store, error) {
	s, err := uc.repo.FindByID(input.StoreID)
	if err != nil {
		return nil, err
	}

	s.UpdateContactInfo(input.ContactPhone)

	if err := uc.repo.Save(s); err != nil {
		return nil, err
	}

	return s, nil
}
