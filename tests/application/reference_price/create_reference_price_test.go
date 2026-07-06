package referenceprice_test

import (
	"testing"

	referencepriceapp "stock-service/internal/application/reference_price"
	"stock-service/internal/domain/reference_price"
)

type inMemoryReferencePriceRepo struct {
	prices map[int64]*referenceprice.ReferencePrice
	nextID int64
}

func newInMemoryReferencePriceRepo() *inMemoryReferencePriceRepo {
	return &inMemoryReferencePriceRepo{
		prices: make(map[int64]*referenceprice.ReferencePrice),
		nextID: 1,
	}
}

func (r *inMemoryReferencePriceRepo) Save(rp *referenceprice.ReferencePrice) error {
	if rp.ID == 0 {
		rp.ID = r.nextID
		r.nextID++
	}
	r.prices[rp.ID] = rp
	return nil
}

func (r *inMemoryReferencePriceRepo) FindByID(id int64) (*referenceprice.ReferencePrice, error) {
	rp, ok := r.prices[id]
	if !ok {
		return nil, nil
	}
	return rp, nil
}

func (r *inMemoryReferencePriceRepo) FindByProductID(productID int32) (*referenceprice.ReferencePrice, error) {
	for _, rp := range r.prices {
		if rp.ProductID == productID {
			return rp, nil
		}
	}
	return nil, nil
}

func (r *inMemoryReferencePriceRepo) FindAll(filter referenceprice.ReferencePriceFilter) ([]*referenceprice.ReferencePrice, int, error) {
	var matched []*referenceprice.ReferencePrice
	for _, rp := range r.prices {
		if filter.ProductID != nil && rp.ProductID != *filter.ProductID {
			continue
		}
		if filter.Source != nil && rp.Source != *filter.Source {
			continue
		}
		matched = append(matched, rp)
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

func (r *inMemoryReferencePriceRepo) Delete(id int64) error {
	delete(r.prices, id)
	return nil
}

func TestCreateReferencePrice_Success(t *testing.T) {
	repo := newInMemoryReferencePriceRepo()
	uc := referencepriceapp.NewCreateReferencePriceUseCase(repo)

	rp, err := uc.Execute(1, 99.99, "supplier")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rp.ID == 0 {
		t.Error("expected ID to be set")
	}
	if rp.ProductID != 1 {
		t.Errorf("expected ProductID %d, got %d", 1, rp.ProductID)
	}
	if rp.Price != 99.99 {
		t.Errorf("expected Price %f, got %f", 99.99, rp.Price)
	}
	if rp.Source != "supplier" {
		t.Errorf("expected Source %q, got %q", "supplier", rp.Source)
	}
}

func TestCreateReferencePrice_ZeroPrice_ReturnsErrInvalidReferencePrice(t *testing.T) {
	repo := newInMemoryReferencePriceRepo()
	uc := referencepriceapp.NewCreateReferencePriceUseCase(repo)

	_, err := uc.Execute(1, 0, "supplier")
	if err != referenceprice.ErrInvalidReferencePrice {
		t.Errorf("expected %v, got %v", referenceprice.ErrInvalidReferencePrice, err)
	}
}

func TestCreateReferencePrice_NegativePrice_ReturnsErrInvalidReferencePrice(t *testing.T) {
	repo := newInMemoryReferencePriceRepo()
	uc := referencepriceapp.NewCreateReferencePriceUseCase(repo)

	_, err := uc.Execute(1, -10, "supplier")
	if err != referenceprice.ErrInvalidReferencePrice {
		t.Errorf("expected %v, got %v", referenceprice.ErrInvalidReferencePrice, err)
	}
}
