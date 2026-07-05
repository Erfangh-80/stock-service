package promotion

import domainpromotion "stock-service/internal/domain/promotion"

type DeactivatePromotionUseCase struct {
	repo domainpromotion.Repository
}

func NewDeactivatePromotionUseCase(repo domainpromotion.Repository) *DeactivatePromotionUseCase {
	return &DeactivatePromotionUseCase{repo: repo}
}

func (uc *DeactivatePromotionUseCase) Execute(id int64) error {
	p, err := uc.repo.FindByID(id)
	if err != nil {
		return err
	}
	p.Deactivate()
	return uc.repo.Save(p)
}
