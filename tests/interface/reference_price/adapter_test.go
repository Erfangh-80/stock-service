package referencepriceinterface_test

import (
	"testing"
	"time"

	app "stock-service/internal/application/reference_price"
	domain "stock-service/internal/domain/reference_price"
	iface "stock-service/internal/interface"
	referencepriceinterface "stock-service/internal/interface/reference_price"
)

type mockCreateReferencePrice struct {
	fn func(int32, float64, string) (*domain.ReferencePrice, error)
}

func (m *mockCreateReferencePrice) Execute(productID int32, price float64, source string) (*domain.ReferencePrice, error) {
	return m.fn(productID, price, source)
}

type mockGetReferencePrice struct {
	fn func(app.GetReferencePriceInput) (*domain.ReferencePrice, error)
}

func (m *mockGetReferencePrice) Execute(input app.GetReferencePriceInput) (*domain.ReferencePrice, error) {
	return m.fn(input)
}

type mockGetByProductReferencePrice struct {
	fn func(app.GetByProductReferencePriceInput) (*domain.ReferencePrice, error)
}

func (m *mockGetByProductReferencePrice) Execute(input app.GetByProductReferencePriceInput) (*domain.ReferencePrice, error) {
	return m.fn(input)
}

type mockListReferencePrices struct {
	fn func(app.ListReferencePricesInput) (*app.ListReferencePricesOutput, error)
}

func (m *mockListReferencePrices) Execute(input app.ListReferencePricesInput) (*app.ListReferencePricesOutput, error) {
	return m.fn(input)
}

type mockDeleteReferencePrice struct {
	fn func(app.DeleteReferencePriceInput) error
}

func (m *mockDeleteReferencePrice) Execute(input app.DeleteReferencePriceInput) error {
	return m.fn(input)
}

type mockValidateReferencePrice struct {
	fn func(app.ValidateReferencePriceInput) (*app.ReferencePriceValidation, error)
}

func (m *mockValidateReferencePrice) Execute(input app.ValidateReferencePriceInput) (*app.ReferencePriceValidation, error) {
	return m.fn(input)
}

func newFullAdapter(
	create *mockCreateReferencePrice,
	get *mockGetReferencePrice,
	getByProduct *mockGetByProductReferencePrice,
	list *mockListReferencePrices,
	del *mockDeleteReferencePrice,
	validate *mockValidateReferencePrice,
) *referencepriceinterface.Adapter {
	return referencepriceinterface.NewAdapter(create, get, getByProduct, list, del, validate)
}

func baseReferencePrice() *domain.ReferencePrice {
	return &domain.ReferencePrice{
		ID: 1, ProductID: 42, Price: 99.99, Source: "supplier",
		CreatedAt: time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
	}
}

func TestAdapter_Create_Success(t *testing.T) {
	rp := baseReferencePrice()
	adapter := newFullAdapter(
		&mockCreateReferencePrice{func(int32, float64, string) (*domain.ReferencePrice, error) { return rp, nil }},
		&mockGetReferencePrice{},
		&mockGetByProductReferencePrice{},
		&mockListReferencePrices{},
		&mockDeleteReferencePrice{},
		&mockValidateReferencePrice{},
	)
	resp, err := adapter.Create(referencepriceinterface.CreateReferencePriceParams{
		ProductID: 42, Price: 99.99, Source: "supplier",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.ProductID != 42 || resp.Price != 99.99 || resp.Source != "supplier" {
		t.Error("unexpected response")
	}
}

func TestAdapter_Create_InvalidInput(t *testing.T) {
	adapter := newFullAdapter(
		&mockCreateReferencePrice{func(int32, float64, string) (*domain.ReferencePrice, error) {
			return nil, domain.ErrInvalidReferencePrice
		}},
		&mockGetReferencePrice{},
		&mockGetByProductReferencePrice{},
		&mockListReferencePrices{},
		&mockDeleteReferencePrice{},
		&mockValidateReferencePrice{},
	)
	_, err := adapter.Create(referencepriceinterface.CreateReferencePriceParams{
		ProductID: 1, Price: 0, Source: "supplier",
	})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

type unknownErr struct{}

func (e unknownErr) Error() string { return "unknown" }

func TestAdapter_Create_UnknownError_ReturnsInternal(t *testing.T) {
	adapter := newFullAdapter(
		&mockCreateReferencePrice{func(int32, float64, string) (*domain.ReferencePrice, error) {
			return nil, unknownErr{}
		}},
		&mockGetReferencePrice{},
		&mockGetByProductReferencePrice{},
		&mockListReferencePrices{},
		&mockDeleteReferencePrice{},
		&mockValidateReferencePrice{},
	)
	_, err := adapter.Create(referencepriceinterface.CreateReferencePriceParams{
		ProductID: 1, Price: 99.99, Source: "supplier",
	})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

func TestAdapter_Get_Success(t *testing.T) {
	rp := baseReferencePrice()
	adapter := newFullAdapter(
		&mockCreateReferencePrice{},
		&mockGetReferencePrice{func(app.GetReferencePriceInput) (*domain.ReferencePrice, error) { return rp, nil }},
		&mockGetByProductReferencePrice{},
		&mockListReferencePrices{},
		&mockDeleteReferencePrice{},
		&mockValidateReferencePrice{},
	)
	resp, err := adapter.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.ProductID != 42 {
		t.Error("unexpected response")
	}
}

func TestAdapter_Get_NotFound(t *testing.T) {
	adapter := newFullAdapter(
		&mockCreateReferencePrice{},
		&mockGetReferencePrice{func(app.GetReferencePriceInput) (*domain.ReferencePrice, error) {
			return nil, domain.ErrReferencePriceNotFound
		}},
		&mockGetByProductReferencePrice{},
		&mockListReferencePrices{},
		&mockDeleteReferencePrice{},
		&mockValidateReferencePrice{},
	)
	_, err := adapter.Get(999)
	if err != iface.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAdapter_GetByProduct_Success(t *testing.T) {
	rp := baseReferencePrice()
	adapter := newFullAdapter(
		&mockCreateReferencePrice{},
		&mockGetReferencePrice{},
		&mockGetByProductReferencePrice{func(app.GetByProductReferencePriceInput) (*domain.ReferencePrice, error) { return rp, nil }},
		&mockListReferencePrices{},
		&mockDeleteReferencePrice{},
		&mockValidateReferencePrice{},
	)
	resp, err := adapter.GetByProduct(42)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ProductID != 42 {
		t.Error("unexpected response")
	}
}

func TestAdapter_List_Success(t *testing.T) {
	output := &app.ListReferencePricesOutput{
		ReferencePrices: []*domain.ReferencePrice{baseReferencePrice()},
		Total: 1, Page: 1, Limit: 20,
	}
	adapter := newFullAdapter(
		&mockCreateReferencePrice{},
		&mockGetReferencePrice{},
		&mockGetByProductReferencePrice{},
		&mockListReferencePrices{func(app.ListReferencePricesInput) (*app.ListReferencePricesOutput, error) { return output, nil }},
		&mockDeleteReferencePrice{},
		&mockValidateReferencePrice{},
	)
	resp, err := adapter.List(referencepriceinterface.ListReferencePricesParams{Limit: 20})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.ReferencePrices) != 1 || resp.Total != 1 || resp.Page != 1 || resp.Limit != 20 {
		t.Error("unexpected list response")
	}
}

func TestAdapter_Delete_Success(t *testing.T) {
	adapter := newFullAdapter(
		&mockCreateReferencePrice{},
		&mockGetReferencePrice{},
		&mockGetByProductReferencePrice{},
		&mockListReferencePrices{},
		&mockDeleteReferencePrice{func(app.DeleteReferencePriceInput) error { return nil }},
		&mockValidateReferencePrice{},
	)
	err := adapter.Delete(1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAdapter_Delete_NotFound(t *testing.T) {
	adapter := newFullAdapter(
		&mockCreateReferencePrice{},
		&mockGetReferencePrice{},
		&mockGetByProductReferencePrice{},
		&mockListReferencePrices{},
		&mockDeleteReferencePrice{func(app.DeleteReferencePriceInput) error { return domain.ErrReferencePriceNotFound }},
		&mockValidateReferencePrice{},
	)
	err := adapter.Delete(999)
	if err != iface.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAdapter_Validate_Success(t *testing.T) {
	validation := &app.ReferencePriceValidation{
		ProductID: 42, ReferencePriceID: 1, ReferencePrice: 99.99,
		Source: "supplier", InventoryCount: 1, BasePrices: []float64{100},
		Comparison: "within_range",
	}
	adapter := newFullAdapter(
		&mockCreateReferencePrice{},
		&mockGetReferencePrice{},
		&mockGetByProductReferencePrice{},
		&mockListReferencePrices{},
		&mockDeleteReferencePrice{},
		&mockValidateReferencePrice{func(app.ValidateReferencePriceInput) (*app.ReferencePriceValidation, error) { return validation, nil }},
	)
	resp, err := adapter.Validate(42)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ProductID != 42 || resp.ReferencePrice != 99.99 || resp.Comparison != "within_range" {
		t.Error("unexpected validation output")
	}
}
