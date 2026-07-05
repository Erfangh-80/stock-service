package inventory_test

import (
	"testing"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
)

func TestSuspendSystemSale_Success(t *testing.T) {
	repo := newInmemoryRepository()
	inv, _ := inventory.NewInventory(1, 1, 42, 100)
	repo.Save(inv)

	uc := appinventory.NewSuspendSystemSaleUseCase(repo)
	result, err := uc.Execute(appinventory.SuspendSystemSaleInput{SaleID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if result.SystemSaleStatus != inventory.SystemSaleStatusSuspended {
		t.Errorf("expected suspended, got %s", result.SystemSaleStatus)
	}
}

func TestSuspendSystemSale_NotFound(t *testing.T) {
	repo := newInmemoryRepository()
	uc := appinventory.NewSuspendSystemSaleUseCase(repo)
	_, err := uc.Execute(appinventory.SuspendSystemSaleInput{SaleID: 999})
	if err != inventory.ErrInventoryNotFound {
		t.Errorf("expected ErrInventoryNotFound, got %v", err)
	}
}
