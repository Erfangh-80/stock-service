package storeallowedcategory

import (
	"stock-service/internal/domain/store_allowed_category"
)

type ApproveCategoryUseCase struct {
	repo storeallowedcategory.Repository
}

func NewApproveCategoryUseCase(repo storeallowedcategory.Repository) *ApproveCategoryUseCase {
	return &ApproveCategoryUseCase{repo: repo}
}

func (uc *ApproveCategoryUseCase) Execute(categoryID int64) error {
	sac, err := uc.repo.FindByID(categoryID)
	if err != nil {
		return err
	}
	sac.Approve()
	return uc.repo.Save(sac)
}
