package store_test

import (
	"testing"

	appstore "stock-service/internal/application/store"
	"stock-service/internal/domain/store"
)

func TestToggleCommissionUseCase_Disable(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	uc := appstore.NewToggleCommissionUseCase(repo)

	s, _ := store.NewStore(1, "My Store")
	repo.Save(s)

	input := appstore.ToggleCommissionInput{StoreID: s.ID}
	result, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.IsCommissionApplicable {
		t.Error("expected IsCommissionApplicable to be false after toggle")
	}
}

func TestToggleCommissionUseCase_Enable(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	uc := appstore.NewToggleCommissionUseCase(repo)

	s, _ := store.NewStore(1, "My Store")
	s.DisableCommission()
	repo.Save(s)

	input := appstore.ToggleCommissionInput{StoreID: s.ID}
	result, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !result.IsCommissionApplicable {
		t.Error("expected IsCommissionApplicable to be true after toggle")
	}
}

func TestToggleCommissionUseCase_StoreNotFound_ReturnsError(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	uc := appstore.NewToggleCommissionUseCase(repo)

	input := appstore.ToggleCommissionInput{StoreID: 999}
	_, err := uc.Execute(input)
	if err != store.ErrStoreNotFound {
		t.Errorf("expected store.ErrStoreNotFound, got %v", err)
	}
}
