package storeallowedcategory

import (
	domainstoreallowedcategory "stock-service/internal/domain/store_allowed_category"
)

type ApproveCategoryInput struct {
	CategoryID int64
}

type ApproveCategoryUseCase struct {
	repo domainstoreallowedcategory.Repository
}

func NewApproveCategoryUseCase(repo domainstoreallowedcategory.Repository) *ApproveCategoryUseCase {
	return &ApproveCategoryUseCase{repo: repo}
}

func (uc *ApproveCategoryUseCase) Execute(input ApproveCategoryInput) error {
	sac, err := uc.repo.FindByID(input.CategoryID)
	if err != nil {
		return err
	}
	if sac == nil {
		return domainstoreallowedcategory.ErrStoreCategoryNotFound
	}
	sac.Approve()
	return uc.repo.Save(sac)
}
