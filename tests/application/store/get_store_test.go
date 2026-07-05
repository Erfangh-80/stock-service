package store_test

import (
	"testing"

	"stock-service/internal/domain/store"
	appstore "stock-service/internal/application/store"
)

func TestGetStoreUseCase_Success(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	createUC := appstore.NewCreateStoreUseCase(repo)
	getUC := appstore.NewGetStoreUseCase(repo)

	created, err := createUC.Execute(appstore.CreateStoreInput{
		UserID: 1, StoreName: "My Store",
	})
	if err != nil {
		t.Fatal(err)
	}

	got, err := getUC.Execute(appstore.GetStoreInput{ID: created.ID})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, got.ID)
	}
}

func TestGetStoreUseCase_NotFound(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	getUC := appstore.NewGetStoreUseCase(repo)

	_, err := getUC.Execute(appstore.GetStoreInput{ID: 999})
	if err != store.ErrStoreNotFound {
		t.Errorf("expected ErrStoreNotFound, got %v", err)
	}
}
