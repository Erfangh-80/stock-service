package salescommission_test

import (
	"testing"

	domainsalescommission "stock-service/internal/domain/sales_commission"
	"stock-service/internal/application/sales_commission"
)

type inMemorySalesCommissionRepo struct {
	commissions map[int64]*domainsalescommission.SalesCommission
	nextID      int64
}

func newInMemorySalesCommissionRepo() *inMemorySalesCommissionRepo {
	return &inMemorySalesCommissionRepo{
		commissions: make(map[int64]*domainsalescommission.SalesCommission),
		nextID:      1,
	}
}

func (r *inMemorySalesCommissionRepo) Save(sc *domainsalescommission.SalesCommission) error {
	if sc.ID == 0 {
		sc.ID = r.nextID
		r.nextID++
	}
	r.commissions[sc.ID] = sc
	return nil
}

func (r *inMemorySalesCommissionRepo) FindByID(id int64) (*domainsalescommission.SalesCommission, error) {
	return r.commissions[id], nil
}

func (r *inMemorySalesCommissionRepo) FindByInventoryID(inventoryID int64) (*domainsalescommission.SalesCommission, error) {
	for _, sc := range r.commissions {
		if sc.InventoryID == inventoryID {
			return sc, nil
		}
	}
	return nil, nil
}

func (r *inMemorySalesCommissionRepo) FindAll(filter domainsalescommission.SalesCommissionFilter) ([]*domainsalescommission.SalesCommission, int, error) {
	var matched []*domainsalescommission.SalesCommission
	for _, sc := range r.commissions {
		if filter.InventoryID != nil && sc.InventoryID != *filter.InventoryID {
			continue
		}
		if filter.SaleModel != nil && string(sc.SaleModel) != *filter.SaleModel {
			continue
		}
		matched = append(matched, sc)
	}
	total := len(matched)
	page, limit := filter.Page, filter.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	start := (page - 1) * limit
	if start >= len(matched) {
		return nil, total, nil
	}
	end := start + limit
	if end > len(matched) {
		end = len(matched)
	}
	return matched[start:end], total, nil
}

func (r *inMemorySalesCommissionRepo) Delete(id int64) error {
	delete(r.commissions, id)
	return nil
}

func TestCreateSalesCommission_Success(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	uc := salescommission.NewCreateSalesCommissionUseCase(repo)

	sc, err := uc.Execute(10, 20, domainsalescommission.SaleModelRetail, 5.0, 100.0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sc.ID == 0 {
		t.Error("expected ID to be set")
	}
	if sc.InventoryID != 10 {
		t.Errorf("expected InventoryID %d, got %d", 10, sc.InventoryID)
	}
	if sc.CategoryCommissionRuleID != 20 {
		t.Errorf("expected CategoryCommissionRuleID %d, got %d", 20, sc.CategoryCommissionRuleID)
	}
	if sc.RatePercent != 5.0 {
		t.Errorf("expected RatePercent %f, got %f", 5.0, sc.RatePercent)
	}
	if sc.MinPrice != 100.0 {
		t.Errorf("expected MinPrice %f, got %f", 100.0, sc.MinPrice)
	}
}

func TestCreateSalesCommission_InvalidRatePercent_ReturnsError(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	uc := salescommission.NewCreateSalesCommissionUseCase(repo)

	_, err := uc.Execute(1, 1, domainsalescommission.SaleModelRetail, 150, 100)
	if err != domainsalescommission.ErrInvalidRatePercent {
		t.Errorf("expected %v, got %v", domainsalescommission.ErrInvalidRatePercent, err)
	}
}

func TestCreateSalesCommission_InvalidMinPrice_ReturnsError(t *testing.T) {
	repo := newInMemorySalesCommissionRepo()
	uc := salescommission.NewCreateSalesCommissionUseCase(repo)

	_, err := uc.Execute(1, 1, domainsalescommission.SaleModelRetail, 5, 0)
	if err != domainsalescommission.ErrInvalidMinPrice {
		t.Errorf("expected %v, got %v", domainsalescommission.ErrInvalidMinPrice, err)
	}
}
