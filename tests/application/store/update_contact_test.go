package store_test

import (
	"testing"

	appstore "stock-service/internal/application/store"
	"stock-service/internal/domain/store"
)

func TestUpdateContactUseCase_Success(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	uc := appstore.NewUpdateContactUseCase(repo)

	s, _ := store.NewStore(1, "My Store")
	repo.Save(s)

	phone := "123-456-7890"
	input := appstore.UpdateContactInput{
		StoreID:      s.ID,
		ContactPhone: &phone,
	}

	result, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ContactPhone == nil {
		t.Fatal("expected ContactPhone to be non-nil")
	}
	if *result.ContactPhone != phone {
		t.Errorf("expected ContactPhone %s, got %s", phone, *result.ContactPhone)
	}
}

func TestUpdateContactUseCase_SetToNil(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	uc := appstore.NewUpdateContactUseCase(repo)

	s, _ := store.NewStore(1, "My Store")
	phone := "old-phone"
	s.ContactPhone = &phone
	repo.Save(s)

	input := appstore.UpdateContactInput{
		StoreID:      s.ID,
		ContactPhone: nil,
	}

	result, err := uc.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ContactPhone != nil {
		t.Errorf("expected ContactPhone to be nil, got %s", *result.ContactPhone)
	}
}

func TestUpdateContactUseCase_StoreNotFound_ReturnsError(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	uc := appstore.NewUpdateContactUseCase(repo)

	input := appstore.UpdateContactInput{
		StoreID:      999,
		ContactPhone: nil,
	}

	_, err := uc.Execute(input)
	if err != store.ErrStoreNotFound {
		t.Errorf("expected store.ErrStoreNotFound, got %v", err)
	}
}
