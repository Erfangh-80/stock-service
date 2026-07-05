package store

import (
	"stock-service/internal/domain/store"
)

type ToggleBulkSaleInput struct {
	StoreID int64
}

type ToggleBulkSaleUseCase interface {
	Execute(input ToggleBulkSaleInput) (*store.Store, error)
}

type toggleBulkSaleInteractor struct {
	repo store.Repository
}

func NewToggleBulkSaleUseCase(repo store.Repository) ToggleBulkSaleUseCase {
	return &toggleBulkSaleInteractor{repo: repo}
}

func (uc *toggleBulkSaleInteractor) Execute(input ToggleBulkSaleInput) (*store.Store, error) {
	s, err := uc.repo.FindByID(input.StoreID)
	if err != nil {
		return nil, err
	}

	if s.IsBulkSaleEnabled {
		s.DisableBulkSale()
	} else {
		s.EnableBulkSale()
	}

	if err := uc.repo.Save(s); err != nil {
		return nil, err
	}

	return s, nil
}
