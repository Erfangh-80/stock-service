package storeallowedcategory

import (
	domainstoreallowedcategory "stock-service/internal/domain/store_allowed_category"
)

type RejectCategoryInput struct {
	CategoryID  int64
	SupportNote string
}

type RejectCategoryUseCase struct {
	repo domainstoreallowedcategory.Repository
}

func NewRejectCategoryUseCase(repo domainstoreallowedcategory.Repository) *RejectCategoryUseCase {
	return &RejectCategoryUseCase{repo: repo}
}

func (uc *RejectCategoryUseCase) Execute(input RejectCategoryInput) error {
	if err := domainstoreallowedcategory.ValidateSupportNote(input.SupportNote); err != nil {
		return err
	}

	sac, err := uc.repo.FindByID(input.CategoryID)
	if err != nil {
		return err
	}
	if sac == nil {
		return domainstoreallowedcategory.ErrStoreCategoryNotFound
	}

	sac.Reject(input.SupportNote)
	return uc.repo.Save(sac)
}
