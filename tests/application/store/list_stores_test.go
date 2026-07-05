package store_test

import (
	"testing"

	appstore "stock-service/internal/application/store"
	"stock-service/internal/domain/store"
)

func TestListStoresUseCase_Success(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	s1, _ := store.NewStore(1, "Store A")
	s1.ID = 1
	repo.Save(s1)
	s2, _ := store.NewStore(2, "Store B")
	s2.ID = 2
	repo.Save(s2)

	uc := appstore.NewListStoresUseCase(repo)

	result, err := uc.Execute(appstore.ListStoresInput{Page: 1, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 2 || len(result.Stores) != 2 {
		t.Errorf("expected 2 stores, got %d total, %d results", result.Total, len(result.Stores))
	}
}

func TestListStoresUseCase_FilterByUser(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	s1, _ := store.NewStore(1, "Store A")
	s1.ID = 1
	repo.Save(s1)
	s2, _ := store.NewStore(1, "Store B")
	s2.ID = 2
	repo.Save(s2)
	s3, _ := store.NewStore(2, "Store C")
	s3.ID = 3
	repo.Save(s3)

	uc := appstore.NewListStoresUseCase(repo)

	uid := int64(1)
	result, err := uc.Execute(appstore.ListStoresInput{UserID: &uid, Page: 1, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 2 {
		t.Errorf("expected 2 stores for user 1, got %d", result.Total)
	}
}

func TestListStoresUseCase_Pagination(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	for i := 0; i < 5; i++ {
		s, _ := store.NewStore(1, "Store")
		s.ID = int64(i + 1)
		repo.Save(s)
	}

	uc := appstore.NewListStoresUseCase(repo)

	result, err := uc.Execute(appstore.ListStoresInput{Page: 1, Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 5 {
		t.Errorf("expected total 5, got %d", result.Total)
	}
	if len(result.Stores) != 2 {
		t.Errorf("expected 2 stores on page 1, got %d", len(result.Stores))
	}

	result2, err := uc.Execute(appstore.ListStoresInput{Page: 3, Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(result2.Stores) != 1 {
		t.Errorf("expected 1 store on page 3, got %d", len(result2.Stores))
	}
}

func TestListStoresUseCase_PageOutOfRange(t *testing.T) {
	t.Parallel()

	repo := newInMemoryStoreRepository()
	s, _ := store.NewStore(1, "Store")
	s.ID = 1
	repo.Save(s)

	uc := appstore.NewListStoresUseCase(repo)

	result, err := uc.Execute(appstore.ListStoresInput{Page: 99, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Total)
	}
	if len(result.Stores) != 0 {
		t.Errorf("expected 0 stores on out-of-range page, got %d", len(result.Stores))
	}
}
