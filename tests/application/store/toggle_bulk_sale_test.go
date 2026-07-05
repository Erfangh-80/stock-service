package store_test

import (
	"testing"

	appstore "stock-service/internal/application/store"
	"stock-service/internal/domain/store"
)

func TestToggleBulkSaleUseCase_Enable(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	uc := appstore.NewToggleBulkSaleUseCase(repo)

	s, _ := store.NewStore(1, "My Store")
	repo.Save(s)

	input := appstore.ToggleBulkSaleInput{StoreID: s.ID}
	result, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !result.IsBulkSaleEnabled {
		t.Error("expected IsBulkSaleEnabled to be true after toggle")
	}
}

func TestToggleBulkSaleUseCase_Disable(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	uc := appstore.NewToggleBulkSaleUseCase(repo)

	s, _ := store.NewStore(1, "My Store")
	s.EnableBulkSale()
	repo.Save(s)

	input := appstore.ToggleBulkSaleInput{StoreID: s.ID}
	result, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.IsBulkSaleEnabled {
		t.Error("expected IsBulkSaleEnabled to be false after toggle")
	}
}

func TestToggleBulkSaleUseCase_StoreNotFound_ReturnsError(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	uc := appstore.NewToggleBulkSaleUseCase(repo)

	input := appstore.ToggleBulkSaleInput{StoreID: 999}
	_, err := uc.Execute(input)
	if err != store.ErrStoreNotFound {
		t.Errorf("expected store.ErrStoreNotFound, got %v", err)
	}
}
