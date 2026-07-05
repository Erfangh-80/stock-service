package inventory_test

import (
	"testing"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
)

func TestReserveQuantity_Success(t *testing.T) {
	repo := newInmemoryRepository()
	inv, _ := inventory.NewInventory(1, 1, 42, 100)
	inv.InstantQty = 50
	repo.Save(inv)

	uc := appinventory.NewReserveQuantityUseCase(repo)
	result, err := uc.Execute(appinventory.ReserveQuantityInput{SaleID: 1, Quantity: 10})
	if err != nil {
		t.Fatal(err)
	}
	if result.InstantQty != 40 {
		t.Errorf("expected qty 40 after reserve, got %d", result.InstantQty)
	}
}

func TestReserveQuantity_InsufficientStock(t *testing.T) {
	repo := newInmemoryRepository()
	inv, _ := inventory.NewInventory(1, 1, 42, 100)
	inv.InstantQty = 5
	repo.Save(inv)

	uc := appinventory.NewReserveQuantityUseCase(repo)
	_, err := uc.Execute(appinventory.ReserveQuantityInput{SaleID: 1, Quantity: 10})
	if err != inventory.ErrInsufficientStock {
		t.Errorf("expected ErrInsufficientStock, got %v", err)
	}
}
