package storewarehouselink

import (
	"stock-service/internal/domain/store_warehouse_link"
)

type ChangeRelationUseCase struct {
	repo storewarehouselink.Repository
}

func NewChangeRelationUseCase(repo storewarehouselink.Repository) *ChangeRelationUseCase {
	return &ChangeRelationUseCase{repo: repo}
}

func (uc *ChangeRelationUseCase) Execute(linkID int64, relationType storewarehouselink.RelationType) error {
	swl, err := uc.repo.FindByID(linkID)
	if err != nil {
		return err
	}
	swl.ChangeRelationType(relationType)
	return uc.repo.Save(swl)
}
