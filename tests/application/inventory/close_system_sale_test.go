package inventory_test

import (
	"testing"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
)

func TestCloseSystemSale_Success(t *testing.T) {
	repo := newInmemoryRepository()
	inv, _ := inventory.NewInventory(1, 1, 42, 100)
	repo.Save(inv)

	uc := appinventory.NewCloseSystemSaleUseCase(repo)
	result, err := uc.Execute(appinventory.CloseSystemSaleInput{SaleID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if result.SystemSaleStatus != inventory.SystemSaleStatusClosed {
		t.Errorf("expected closed, got %s", result.SystemSaleStatus)
	}
}

func TestCloseSystemSale_AlreadyClosed_ReturnsError(t *testing.T) {
	repo := newInmemoryRepository()
	inv, _ := inventory.NewInventory(1, 1, 42, 100)
	inv.CloseSystemSale()
	repo.Save(inv)

	uc := appinventory.NewCloseSystemSaleUseCase(repo)
	_, err := uc.Execute(appinventory.CloseSystemSaleInput{SaleID: 1})
	if err != inventory.ErrSystemSaleStatusTransition {
		t.Errorf("expected ErrSystemSaleStatusTransition, got %v", err)
	}
}
