package store_test

import (
	"testing"

	appstore "stock-service/internal/application/store"
	"stock-service/internal/domain/store"
)

func TestUpdateStoreNameUseCase_Success(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	s, _ := store.NewStore(1, "Original")
	s.ID = 1
	repo.Save(s)

	uc := appstore.NewUpdateStoreNameUseCase(repo)

	result, err := uc.Execute(appstore.UpdateStoreNameInput{StoreID: 1, Name: "Updated"})
	if err != nil {
		t.Fatal(err)
	}
	if result.StoreName != "Updated" {
		t.Errorf("expected 'Updated', got %q", result.StoreName)
	}

	// verify persisted
	saved, _ := repo.FindByID(1)
	if saved.StoreName != "Updated" {
		t.Errorf("persisted name should be 'Updated', got %q", saved.StoreName)
	}
}

func TestUpdateStoreNameUseCase_NotFound(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	uc := appstore.NewUpdateStoreNameUseCase(repo)

	_, err := uc.Execute(appstore.UpdateStoreNameInput{StoreID: 999, Name: "Nope"})
	if err != store.ErrStoreNotFound {
		t.Errorf("expected ErrStoreNotFound, got %v", err)
	}
}

func TestUpdateStoreNameUseCase_EmptyName(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	s, _ := store.NewStore(1, "Original")
	s.ID = 1
	repo.Save(s)

	uc := appstore.NewUpdateStoreNameUseCase(repo)

	_, err := uc.Execute(appstore.UpdateStoreNameInput{StoreID: 1, Name: ""})
	if err != store.ErrStoreNameRequired {
		t.Errorf("expected ErrStoreNameRequired, got %v", err)
	}
}
