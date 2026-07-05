package inventory_test

import (
	"testing"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
)

func TestReleaseQuantity_Success(t *testing.T) {
	repo := newInmemoryRepository()
	inv, _ := inventory.NewInventory(1, 1, 42, 100)
	inv.InstantQty = 50
	repo.Save(inv)

	uc := appinventory.NewReleaseQuantityUseCase(repo)
	result, err := uc.Execute(appinventory.ReleaseQuantityInput{SaleID: 1, Quantity: 10})
	if err != nil {
		t.Fatal(err)
	}
	if result.InstantQty != 60 {
		t.Errorf("expected qty 60 after release, got %d", result.InstantQty)
	}
}

func TestReleaseQuantity_NotFound(t *testing.T) {
	repo := newInmemoryRepository()
	uc := appinventory.NewReleaseQuantityUseCase(repo)
	_, err := uc.Execute(appinventory.ReleaseQuantityInput{SaleID: 999, Quantity: 10})
	if err != inventory.ErrInventoryNotFound {
		t.Errorf("expected ErrInventoryNotFound, got %v", err)
	}
}
