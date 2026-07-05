package store_test

import (
	"testing"

	appstore "stock-service/internal/application/store"
	"stock-service/internal/domain/store"
)

func TestCreateStoreUseCase_Success(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	uc := appstore.NewCreateStoreUseCase(repo)

	input := appstore.CreateStoreInput{
		UserID:    1,
		StoreName: "My Store",
	}

	s, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s.UserID != 1 {
		t.Errorf("expected UserID 1, got %d", s.UserID)
	}
	if s.StoreName != "My Store" {
		t.Errorf("expected StoreName 'My Store', got %s", s.StoreName)
	}
	if s.Status != store.StoreStatusActive {
		t.Errorf("expected Status %q, got %q", store.StoreStatusActive, s.Status)
	}
	if !s.IsCommissionApplicable {
		t.Error("expected IsCommissionApplicable to be true")
	}
	if s.IsBulkSaleEnabled {
		t.Error("expected IsBulkSaleEnabled to be false")
	}
	if s.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestCreateStoreUseCase_EmptyName_ReturnsError(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	uc := appstore.NewCreateStoreUseCase(repo)

	input := appstore.CreateStoreInput{
		UserID:    1,
		StoreName: "",
	}

	_, err := uc.Execute(input)
	if err != store.ErrStoreNameRequired {
		t.Errorf("expected ErrStoreNameRequired, got %v", err)
	}
}

func TestCreateStoreUseCase_NameTooLong_ReturnsError(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	uc := appstore.NewCreateStoreUseCase(repo)

	input := appstore.CreateStoreInput{
		UserID:    1,
		StoreName: string(make([]byte, 256)),
	}

	_, err := uc.Execute(input)
	if err != store.ErrStoreNameTooLong {
		t.Errorf("expected ErrStoreNameTooLong, got %v", err)
	}
}
