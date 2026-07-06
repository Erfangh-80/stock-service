package salescommission_test

import (
	"testing"

	domainsalescommission "stock-service/internal/domain/sales_commission"
	"stock-service/internal/application/sales_commission"
)

func TestUpdateMinQty_Success(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	uc := salescommission.NewUpdateMinQtyUseCase(repo)

	sc, _ := domainsalescommission.NewSalesCommission(1, 1, domainsalescommission.SaleModelRetail, 5, 100)
	repo.Save(sc)

	err := uc.Execute(sc.ID, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	saved, _ := repo.FindByID(sc.ID)
	if saved.MinQty == nil {
		t.Fatal("expected MinQty to be set, got nil")
	}
	if *saved.MinQty != 10 {
		t.Errorf("expected MinQty %d, got %d", 10, *saved.MinQty)
	}
}

func TestUpdateMinQty_NotFound_ReturnsError(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	uc := salescommission.NewUpdateMinQtyUseCase(repo)

	err := uc.Execute(999, 10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateMinQty_NegativeQty_ReturnsError(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	uc := salescommission.NewUpdateMinQtyUseCase(repo)

	sc, _ := domainsalescommission.NewSalesCommission(1, 1, domainsalescommission.SaleModelRetail, 5, 100)
	repo.Save(sc)

	err := uc.Execute(sc.ID, -1)
	if err != domainsalescommission.ErrInvalidMinQty {
		t.Errorf("expected %v, got %v", domainsalescommission.ErrInvalidMinQty, err)
	}
}
