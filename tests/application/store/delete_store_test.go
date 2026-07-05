package store_test

import (
	"testing"

	appstore "stock-service/internal/application/store"
	"stock-service/internal/domain/store"
)

func TestDeleteStoreUseCase_Success(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	s, _ := store.NewStore(1, "Store")
	s.ID = 1
	repo.Save(s)

	uc := appstore.NewDeleteStoreUseCase(repo)

	err := uc.Execute(appstore.DeleteStoreInput{ID: 1})
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.FindByID(1)
	if err != store.ErrStoreNotFound {
		t.Error("expected store to be deleted")
	}
}

func TestDeleteStoreUseCase_NotFound(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	uc := appstore.NewDeleteStoreUseCase(repo)

	err := uc.Execute(appstore.DeleteStoreInput{ID: 999})
	if err != store.ErrStoreNotFound {
		t.Errorf("expected ErrStoreNotFound, got %v", err)
	}
}
