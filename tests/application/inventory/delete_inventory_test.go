package inventory_test

import (
	"testing"

	appinventory "stock-service/internal/application/inventory"
	"stock-service/internal/domain/inventory"
)

func TestDeleteInventory_Success(t *testing.T) {
	repo := newInmemoryRepository()
	inv, _ := inventory.NewInventory(1, 1, 42, 100)
	repo.Save(inv)

	uc := appinventory.NewDeleteInventoryUseCase(repo)
	err := uc.Execute(appinventory.DeleteInventoryInput{SaleID: 1})
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.FindByID(1)
	if err != inventory.ErrInventoryNotFound {
		t.Errorf("expected ErrInventoryNotFound after delete, got %v", err)
	}
}

func TestDeleteInventory_NotFound(t *testing.T) {
	repo := newInmemoryRepository()
	uc := appinventory.NewDeleteInventoryUseCase(repo)
	err := uc.Execute(appinventory.DeleteInventoryInput{SaleID: 999})
	if err != inventory.ErrInventoryNotFound {
		t.Errorf("expected ErrInventoryNotFound, got %v", err)
	}
}
