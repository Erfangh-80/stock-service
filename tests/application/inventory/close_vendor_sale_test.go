package inventory_test

import (
	"testing"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
)

func TestCloseVendorSale_Success(t *testing.T) {
	repo := newInmemoryRepository()
	inv, _ := inventory.NewInventory(1, 1, 42, 100)
	repo.Save(inv)

	uc := appinventory.NewCloseVendorSaleUseCase(repo)
	result, err := uc.Execute(appinventory.CloseVendorSaleInput{SaleID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if result.VendorSaleStatus != inventory.VendorSaleStatusClosed {
		t.Errorf("expected closed, got %s", result.VendorSaleStatus)
	}
}

func TestCloseVendorSale_AlreadyClosed_ReturnsError(t *testing.T) {
	repo := newInmemoryRepository()
	inv, _ := inventory.NewInventory(1, 1, 42, 100)
	inv.CloseVendorSale()
	repo.Save(inv)

	uc := appinventory.NewCloseVendorSaleUseCase(repo)
	_, err := uc.Execute(appinventory.CloseVendorSaleInput{SaleID: 1})
	if err != inventory.ErrVendorSaleStatusTransition {
		t.Errorf("expected ErrVendorSaleStatusTransition, got %v", err)
	}
}
