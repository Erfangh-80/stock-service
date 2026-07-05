package salescommission_test

import (
	"errors"
	"testing"

	domainsalescommission "stock-service/internal/domain/sales_commission"
	"stock-service/internal/application/sales_commission"
)

type updateMinQtyInMemoryRepo struct {
	commissions map[int64]*domainsalescommission.SalesCommission
	nextID      int64
}

func newUpdateMinQtyInMemoryRepo() *updateMinQtyInMemoryRepo {
	return &updateMinQtyInMemoryRepo{
		commissions: make(map[int64]*domainsalescommission.SalesCommission),
		nextID:      1,
	}
}

func (r *updateMinQtyInMemoryRepo) Save(sc *domainsalescommission.SalesCommission) error {
	if sc.ID == 0 {
		sc.ID = r.nextID
		r.nextID++
	}
	r.commissions[sc.ID] = sc
	return nil
}

func (r *updateMinQtyInMemoryRepo) FindByID(id int64) (*domainsalescommission.SalesCommission, error) {
	sc, ok := r.commissions[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return sc, nil
}

func (r *updateMinQtyInMemoryRepo) Delete(id int64) error {
	delete(r.commissions, id)
	return nil
}

func TestUpdateMinQty_Success(t *testing.T) {
	repo := newUpdateMinQtyInMemoryRepo()
	uc := salescommission.NewUpdateMinQtyUseCase(repo)

	sc, _ := domainsalescommission.NewSalesCommission(1, 1, domainsalescommission.SaleModelRetail, 5, 100)
	repo.Save(sc)

	err := uc.Execute(sc.ID, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	saved, _ := repo.FindByID(sc.ID)
	if saved.MinQty == nil {
		t.Fatal("expected MinQty to be set, got nil")
	}
	if *saved.MinQty != 10 {
		t.Errorf("expected MinQty %d, got %d", 10, *saved.MinQty)
	}
}

func TestUpdateMinQty_NotFound_ReturnsError(t *testing.T) {
	repo := newUpdateMinQtyInMemoryRepo()
	uc := salescommission.NewUpdateMinQtyUseCase(repo)

	err := uc.Execute(999, 10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateMinQty_NegativeQty_ReturnsError(t *testing.T) {
	repo := newUpdateMinQtyInMemoryRepo()
	uc := salescommission.NewUpdateMinQtyUseCase(repo)

	sc, _ := domainsalescommission.NewSalesCommission(1, 1, domainsalescommission.SaleModelRetail, 5, 100)
	repo.Save(sc)

	err := uc.Execute(sc.ID, -1)
	if err != domainsalescommission.ErrInvalidMinQty {
		t.Errorf("expected %v, got %v", domainsalescommission.ErrInvalidMinQty, err)
	}
}
