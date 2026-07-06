package salescommission_test

import (
	"testing"

	domainsalescommission "stock-service/internal/domain/sales_commission"
	"stock-service/internal/application/sales_commission"
)

func TestUpdateMaxPrice_Success(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	uc := salescommission.NewUpdateMaxPriceUseCase(repo)

	sc, _ := domainsalescommission.NewSalesCommission(1, 1, domainsalescommission.SaleModelRetail, 5, 100)
	repo.Save(sc)

	err := uc.Execute(sc.ID, 200)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	saved, _ := repo.FindByID(sc.ID)
	if saved.MaxPrice == nil {
		t.Fatal("expected MaxPrice to be set, got nil")
	}
	if *saved.MaxPrice != 200 {
		t.Errorf("expected MaxPrice %f, got %f", 200.0, *saved.MaxPrice)
	}
}

func TestUpdateMaxPrice_NotFound_ReturnsError(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	uc := salescommission.NewUpdateMaxPriceUseCase(repo)

	err := uc.Execute(999, 200)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateMaxPrice_InvalidMaxPrice_ReturnsError(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	uc := salescommission.NewUpdateMaxPriceUseCase(repo)

	sc, _ := domainsalescommission.NewSalesCommission(1, 1, domainsalescommission.SaleModelRetail, 5, 100)
	repo.Save(sc)

	err := uc.Execute(sc.ID, 50)
	if err != domainsalescommission.ErrInvalidMaxPrice {
		t.Errorf("expected %v, got %v", domainsalescommission.ErrInvalidMaxPrice, err)
	}
}
