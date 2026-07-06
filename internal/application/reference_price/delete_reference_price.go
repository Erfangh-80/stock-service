package referenceprice

import domainreferenceprice "stock-service/internal/domain/reference_price"

type DeleteReferencePriceInput struct {
	ID int64
}

type DeleteReferencePriceUseCase struct {
	repo domainreferenceprice.Repository
}

func NewDeleteReferencePriceUseCase(repo domainreferenceprice.Repository) *DeleteReferencePriceUseCase {
	return &DeleteReferencePriceUseCase{repo: repo}
}

func (uc *DeleteReferencePriceUseCase) Execute(input DeleteReferencePriceInput) error {
	rp, err := uc.repo.FindByID(input.ID)
	if err != nil {
		return err
	}
	if rp == nil {
		return domainreferenceprice.ErrReferencePriceNotFound
	}
	return uc.repo.Delete(input.ID)
}
