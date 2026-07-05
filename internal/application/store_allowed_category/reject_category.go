package storeallowedcategory

import (
	"stock-service/internal/domain/store_allowed_category"
)

type RejectCategoryUseCase struct {
	repo storeallowedcategory.Repository
}

func NewRejectCategoryUseCase(repo storeallowedcategory.Repository) *RejectCategoryUseCase {
	return &RejectCategoryUseCase{repo: repo}
}

func (uc *RejectCategoryUseCase) Execute(categoryID int64) error {
	sac, err := uc.repo.FindByID(categoryID)
	if err != nil {
		return err
	}
	sac.Reject()
	return uc.repo.Save(sac)
}
