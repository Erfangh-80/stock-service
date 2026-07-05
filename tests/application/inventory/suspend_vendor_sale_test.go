package inventory_test

import (
	"testing"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
)

func TestSuspendVendorSale_Success(t *testing.T) {
	repo := newInmemoryRepository()
	inv, _ := inventory.NewInventory(1, 1, 42, 100)
	repo.Save(inv)

	uc := appinventory.NewSuspendVendorSaleUseCase(repo)
	result, err := uc.Execute(appinventory.SuspendVendorSaleInput{SaleID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if result.VendorSaleStatus != inventory.VendorSaleStatusSuspended {
		t.Errorf("expected suspended, got %s", result.VendorSaleStatus)
	}
}

func TestSuspendVendorSale_NotFound(t *testing.T) {
	repo := newInmemoryRepository()
	uc := appinventory.NewSuspendVendorSaleUseCase(repo)
	_, err := uc.Execute(appinventory.SuspendVendorSaleInput{SaleID: 999})
	if err != inventory.ErrInventoryNotFound {
		t.Errorf("expected ErrInventoryNotFound, got %v", err)
	}
}
