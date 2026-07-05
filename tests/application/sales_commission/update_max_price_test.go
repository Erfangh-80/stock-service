package salescommission_test

import (
	"errors"
	"testing"

	domainsalescommission "stock-service/internal/domain/sales_commission"
	"stock-service/internal/application/sales_commission"
)

type updateMaxPriceInMemoryRepo struct {
	commissions map[int64]*domainsalescommission.SalesCommission
	nextID      int64
}

func newUpdateMaxPriceInMemoryRepo() *updateMaxPriceInMemoryRepo {
	return &updateMaxPriceInMemoryRepo{
		commissions: make(map[int64]*domainsalescommission.SalesCommission),
		nextID:      1,
	}
}

func (r *updateMaxPriceInMemoryRepo) Save(sc *domainsalescommission.SalesCommission) error {
	if sc.ID == 0 {
		sc.ID = r.nextID
		r.nextID++
	}
	r.commissions[sc.ID] = sc
	return nil
}

func (r *updateMaxPriceInMemoryRepo) FindByID(id int64) (*domainsalescommission.SalesCommission, error) {
	sc, ok := r.commissions[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return sc, nil
}

func (r *updateMaxPriceInMemoryRepo) Delete(id int64) error {
	delete(r.commissions, id)
	return nil
}

func TestUpdateMaxPrice_Success(t *testing.T) {
	repo := newUpdateMaxPriceInMemoryRepo()
	uc := salescommission.NewUpdateMaxPriceUseCase(repo)

	sc, _ := domainsalescommission.NewSalesCommission(1, 1, domainsalescommission.SaleModelRetail, 5, 100)
	repo.Save(sc)

	err := uc.Execute(sc.ID, 200)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	saved, _ := repo.FindByID(sc.ID)
	if saved.MaxPrice == nil {
		t.Fatal("expected MaxPrice to be set, got nil")
	}
	if *saved.MaxPrice != 200 {
		t.Errorf("expected MaxPrice %f, got %f", 200.0, *saved.MaxPrice)
	}
}

func TestUpdateMaxPrice_NotFound_ReturnsError(t *testing.T) {
	repo := newUpdateMaxPriceInMemoryRepo()
	uc := salescommission.NewUpdateMaxPriceUseCase(repo)

	err := uc.Execute(999, 200)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateMaxPrice_InvalidMaxPrice_ReturnsError(t *testing.T) {
	repo := newUpdateMaxPriceInMemoryRepo()
	uc := salescommission.NewUpdateMaxPriceUseCase(repo)

	sc, _ := domainsalescommission.NewSalesCommission(1, 1, domainsalescommission.SaleModelRetail, 5, 100)
	repo.Save(sc)

	err := uc.Execute(sc.ID, 50)
	if err != domainsalescommission.ErrInvalidMaxPrice {
		t.Errorf("expected %v, got %v", domainsalescommission.ErrInvalidMaxPrice, err)
	}
}
