package salescommission_test

import (
	"testing"

	domainsalescommission "stock-service/internal/domain/sales_commission"
	"stock-service/internal/application/sales_commission"
)

func TestGetSalesCommission_Success(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	sc, _ := domainsalescommission.NewSalesCommission(10, 20, domainsalescommission.SaleModelRetail, 5, 100)
	repo.Save(sc)

	uc := salescommission.NewGetSalesCommissionUseCase(repo)
	result, err := uc.Execute(salescommission.GetSalesCommissionInput{ID: sc.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.InventoryID != 10 {
		t.Errorf("expected InventoryID 10, got %d", result.InventoryID)
	}
}

func TestGetSalesCommission_NotFound(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	uc := salescommission.NewGetSalesCommissionUseCase(repo)

	_, err := uc.Execute(salescommission.GetSalesCommissionInput{ID: 999})
	if err != domainsalescommission.ErrCommissionNotFound {
		t.Errorf("expected ErrCommissionNotFound, got %v", err)
	}
}

func TestGetByInventorySalesCommission_Success(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	sc, _ := domainsalescommission.NewSalesCommission(10, 20, domainsalescommission.SaleModelRetail, 5, 100)
	repo.Save(sc)

	uc := salescommission.NewGetByInventorySalesCommissionUseCase(repo)
	result, err := uc.Execute(salescommission.GetByInventorySalesCommissionInput{InventoryID: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.InventoryID != 10 {
		t.Errorf("expected InventoryID 10, got %d", result.InventoryID)
	}
}

func TestGetByInventorySalesCommission_NotFound(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	uc := salescommission.NewGetByInventorySalesCommissionUseCase(repo)

	_, err := uc.Execute(salescommission.GetByInventorySalesCommissionInput{InventoryID: 999})
	if err != domainsalescommission.ErrCommissionNotFound {
		t.Errorf("expected ErrCommissionNotFound, got %v", err)
	}
}

func TestDeleteSalesCommission_Success(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	sc, _ := domainsalescommission.NewSalesCommission(10, 20, domainsalescommission.SaleModelRetail, 5, 100)
	repo.Save(sc)

	uc := salescommission.NewDeleteSalesCommissionUseCase(repo)
	err := uc.Execute(salescommission.DeleteSalesCommissionInput{ID: sc.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	saved, _ := repo.FindByID(sc.ID)
	if saved != nil {
		t.Error("expected commission to be deleted")
	}
}

func TestDeleteSalesCommission_NotFound(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	uc := salescommission.NewDeleteSalesCommissionUseCase(repo)

	err := uc.Execute(salescommission.DeleteSalesCommissionInput{ID: 999})
	if err != domainsalescommission.ErrCommissionNotFound {
		t.Errorf("expected ErrCommissionNotFound, got %v", err)
	}
}

func TestListSalesCommissions_Empty(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	uc := salescommission.NewListSalesCommissionsUseCase(repo)

	result, err := uc.Execute(salescommission.ListSalesCommissionsInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Commissions) != 0 {
		t.Errorf("expected 0 items, got %d", len(result.Commissions))
	}
}

func TestListSalesCommissions_All(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	repo.Save(&domainsalescommission.SalesCommission{InventoryID: 1, SaleModel: "retail", RatePercent: 5, MinPrice: 100})
	repo.Save(&domainsalescommission.SalesCommission{InventoryID: 2, SaleModel: "retail", RatePercent: 10, MinPrice: 200})

	uc := salescommission.NewListSalesCommissionsUseCase(repo)
	result, err := uc.Execute(salescommission.ListSalesCommissionsInput{Limit: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Commissions) != 2 {
		t.Errorf("expected 2 items, got %d", len(result.Commissions))
	}
}
