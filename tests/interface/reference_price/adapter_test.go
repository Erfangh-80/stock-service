package referencepriceinterface_test

import (
	"testing"
	"time"

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

func baseReferencePrice() *domain.ReferencePrice {
	return &domain.ReferencePrice{
		ID: 1, ProductID: 42, Price: 99.99, Source: "supplier",
		CreatedAt: time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
	}
}

func TestAdapter_Create_Success(t *testing.T) {
	rp := baseReferencePrice()
	adapter := referencepriceinterface.NewAdapter(
		&mockCreateReferencePrice{func(productID int32, price float64, source string) (*domain.ReferencePrice, error) {
			return rp, nil
		}},
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
	adapter := referencepriceinterface.NewAdapter(
		&mockCreateReferencePrice{func(productID int32, price float64, source string) (*domain.ReferencePrice, error) {
			return nil, domain.ErrInvalidReferencePrice
		}},
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
	adapter := referencepriceinterface.NewAdapter(
		&mockCreateReferencePrice{func(productID int32, price float64, source string) (*domain.ReferencePrice, error) {
			return nil, unknownErr{}
		}},
	)
	_, err := adapter.Create(referencepriceinterface.CreateReferencePriceParams{
		ProductID: 1, Price: 99.99, Source: "supplier",
	})
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}
