package referenceprice

import (
	"math"

	domainreferenceprice "stock-service/internal/domain/reference_price"
	inventorydomain "stock-service/internal/domain/inventory"
)

type ValidateReferencePriceInput struct {
	ProductID int32
}

type ReferencePriceValidation struct {
	ProductID        int32
	ReferencePriceID int64
	ReferencePrice   float64
	Source           string
	InventoryCount   int
	BasePrices       []float64
	Comparison       string
}

type ValidateReferencePriceUseCase struct {
	refPriceRepo domainreferenceprice.Repository
	invRepo      inventorydomain.Repository
}

func NewValidateReferencePriceUseCase(refPriceRepo domainreferenceprice.Repository, invRepo inventorydomain.Repository) *ValidateReferencePriceUseCase {
	return &ValidateReferencePriceUseCase{refPriceRepo: refPriceRepo, invRepo: invRepo}
}

func (uc *ValidateReferencePriceUseCase) Execute(input ValidateReferencePriceInput) (*ReferencePriceValidation, error) {
	rp, err := uc.refPriceRepo.FindByProductID(input.ProductID)
	if err != nil {
		return nil, err
	}
	if rp == nil {
		return nil, domainreferenceprice.ErrReferencePriceNotFound
	}

	allInventory, err := uc.invRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var basePrices []float64
	for _, inv := range allInventory {
		if inv.ProductID == input.ProductID {
			basePrices = append(basePrices, inv.BasePrice)
		}
	}

	result := &ReferencePriceValidation{
		ProductID:        input.ProductID,
		ReferencePriceID: rp.ID,
		ReferencePrice:   rp.Price,
		Source:           rp.Source,
		InventoryCount:   len(basePrices),
		BasePrices:       basePrices,
	}

	if len(basePrices) == 0 {
		result.Comparison = "no_inventory"
	} else {
		avgBasePrice := 0.0
		for _, bp := range basePrices {
			avgBasePrice += bp
		}
		avgBasePrice /= float64(len(basePrices))

		diff := rp.Price - avgBasePrice
		diffPct := (diff / avgBasePrice) * 100

		diffStr := "within_range"
		absDiffPct := math.Abs(diffPct)
		if absDiffPct > 50 {
			diffStr = "far_out_of_range"
		} else if absDiffPct > 20 {
			diffStr = "moderately_out_of_range"
		}

		result.Comparison = diffStr
	}

	return result, nil
}
