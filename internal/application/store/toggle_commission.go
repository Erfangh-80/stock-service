package store

import (
	"stock-service/internal/domain/store"
)

type ToggleCommissionInput struct {
	StoreID int64
}

type ToggleCommissionUseCase interface {
	Execute(input ToggleCommissionInput) (*store.Store, error)
}

type toggleCommissionInteractor struct {
	repo store.Repository
}

func NewToggleCommissionUseCase(repo store.Repository) ToggleCommissionUseCase {
	return &toggleCommissionInteractor{repo: repo}
}

func (uc *toggleCommissionInteractor) Execute(input ToggleCommissionInput) (*store.Store, error) {
	s, err := uc.repo.FindByID(input.StoreID)
	if err != nil {
		return nil, err
	}

	if s.IsCommissionApplicable {
		s.DisableCommission()
	} else {
		s.EnableCommission()
	}

	if err := uc.repo.Save(s); err != nil {
		return nil, err
	}

	return s, nil
}
