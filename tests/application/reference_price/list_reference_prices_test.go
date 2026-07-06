package referenceprice_test

import (
	"testing"

	referencepriceapp "stock-service/internal/application/reference_price"
	"stock-service/internal/domain/reference_price"
)

func TestListReferencePrices_Empty(t *testing.T) {
	repo := newInMemoryReferencePriceRepo()
	uc := referencepriceapp.NewListReferencePricesUseCase(repo)

	result, err := uc.Execute(referencepriceapp.ListReferencePricesInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.ReferencePrices) != 0 {
		t.Errorf("expected 0 items, got %d", len(result.ReferencePrices))
	}
	if result.Total != 0 {
		t.Errorf("expected total 0, got %d", result.Total)
	}
}

func TestListReferencePrices_All(t *testing.T) {
	repo := newInMemoryReferencePriceRepo()
	repo.Save(&referenceprice.ReferencePrice{ProductID: 1, Price: 10, Source: "a"})
	repo.Save(&referenceprice.ReferencePrice{ProductID: 2, Price: 20, Source: "b"})
	repo.Save(&referenceprice.ReferencePrice{ProductID: 3, Price: 30, Source: "a"})

	uc := referencepriceapp.NewListReferencePricesUseCase(repo)

	result, err := uc.Execute(referencepriceapp.ListReferencePricesInput{Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.ReferencePrices) != 3 {
		t.Errorf("expected 3 items, got %d", len(result.ReferencePrices))
	}
	if result.Total != 3 {
		t.Errorf("expected total 3, got %d", result.Total)
	}
}

func TestListReferencePrices_FilterBySource(t *testing.T) {
	repo := newInMemoryReferencePriceRepo()
	repo.Save(&referenceprice.ReferencePrice{ProductID: 1, Price: 10, Source: "a"})
	repo.Save(&referenceprice.ReferencePrice{ProductID: 2, Price: 20, Source: "b"})

	source := "a"
	uc := referencepriceapp.NewListReferencePricesUseCase(repo)

	result, err := uc.Execute(referencepriceapp.ListReferencePricesInput{Source: &source})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.ReferencePrices) != 1 {
		t.Errorf("expected 1 item, got %d", len(result.ReferencePrices))
	}
}

func TestListReferencePrices_Pagination(t *testing.T) {
	repo := newInMemoryReferencePriceRepo()
	for i := int32(1); i <= 10; i++ {
		repo.Save(&referenceprice.ReferencePrice{ProductID: i, Price: float64(i) * 10, Source: "x"})
	}

	uc := referencepriceapp.NewListReferencePricesUseCase(repo)

	result, err := uc.Execute(referencepriceapp.ListReferencePricesInput{Page: 2, Limit: 3})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.ReferencePrices) != 3 {
		t.Errorf("expected 3 items on page 2, got %d", len(result.ReferencePrices))
	}
	if result.Total != 10 {
		t.Errorf("expected total 10, got %d", result.Total)
	}
	if result.Page != 2 {
		t.Errorf("expected page 2, got %d", result.Page)
	}
	if result.Limit != 3 {
		t.Errorf("expected limit 3, got %d", result.Limit)
	}
}
