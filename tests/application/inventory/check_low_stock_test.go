package inventory_test

import (
	"testing"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
)

func TestCheckLowStock_IsLow(t *testing.T) {
	repo := newInmemoryRepository()
	inv, _ := inventory.NewInventory(1, 1, 42, 100)
	inv.InstantQty = 3
	repo.Save(inv)

	uc := appinventory.NewCheckLowStockUseCase(repo)
	result, err := uc.Execute(appinventory.CheckLowStockInput{SaleID: 1, Threshold: 5})
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsLow {
		t.Error("expected IsLow true")
	}
	if result.CurrentQty != 3 {
		t.Errorf("expected qty 3, got %d", result.CurrentQty)
	}
}

func TestCheckLowStock_NotLow(t *testing.T) {
	repo := newInmemoryRepository()
	inv, _ := inventory.NewInventory(1, 1, 42, 100)
	inv.InstantQty = 20
	repo.Save(inv)

	uc := appinventory.NewCheckLowStockUseCase(repo)
	result, err := uc.Execute(appinventory.CheckLowStockInput{SaleID: 1, Threshold: 5})
	if err != nil {
		t.Fatal(err)
	}
	if result.IsLow {
		t.Error("expected IsLow false")
	}
}
