package promotion

import domainpromotion "stock-service/internal/domain/promotion"

type ActivatePromotionUseCase struct {
	repo domainpromotion.Repository
}

func NewActivatePromotionUseCase(repo domainpromotion.Repository) *ActivatePromotionUseCase {
	return &ActivatePromotionUseCase{repo: repo}
}

func (uc *ActivatePromotionUseCase) Execute(id int64) error {
	p, err := uc.repo.FindByID(id)
	if err != nil {
		return err
	}
	p.Activate()
	return uc.repo.Save(p)
}
