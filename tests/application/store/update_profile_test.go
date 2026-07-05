package store_test

import (
	"testing"

	appstore "stock-service/internal/application/store"
	"stock-service/internal/domain/store"
)

func TestUpdateStoreProfileUseCase_AddressID(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	s, _ := store.NewStore(1, "Store")
	s.ID = 1
	repo.Save(s)

	addr := int64(42)
	uc := appstore.NewUpdateStoreProfileUseCase(repo)

	result, err := uc.Execute(appstore.UpdateStoreProfileInput{StoreID: 1, AddressID: &addr})
	if err != nil {
		t.Fatal(err)
	}
	if result.AddressID == nil || *result.AddressID != 42 {
		t.Errorf("expected address_id 42, got %v", result.AddressID)
	}
}

func TestUpdateStoreProfileUseCase_MediaAssets(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	s, _ := store.NewStore(1, "Store")
	s.ID = 1
	repo.Save(s)

	assets := map[string]any{"logo": "logo.png"}
	uc := appstore.NewUpdateStoreProfileUseCase(repo)

	result, err := uc.Execute(appstore.UpdateStoreProfileInput{StoreID: 1, MediaAssets: assets})
	if err != nil {
		t.Fatal(err)
	}
	if result.MediaAssets == nil || result.MediaAssets["logo"] != "logo.png" {
		t.Errorf("expected media_assets to contain logo")
	}
}

func TestUpdateStoreProfileUseCase_NotFound(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	uc := appstore.NewUpdateStoreProfileUseCase(repo)

	_, err := uc.Execute(appstore.UpdateStoreProfileInput{StoreID: 999})
	if err != store.ErrStoreNotFound {
		t.Errorf("expected ErrStoreNotFound, got %v", err)
	}
}
