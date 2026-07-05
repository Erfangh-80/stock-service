package referenceprice

import (
	"testing"

	"stock-service/internal/domain/reference_price"
)

func TestNewReferencePrice_ValidInputs_Succeeds(t *testing.T) {
	rp, err := referenceprice.NewReferencePrice(1, 99.99, "supplier")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
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

func TestNewReferencePrice_ZeroPrice_ReturnsErrInvalidReferencePrice(t *testing.T) {
	_, err := referenceprice.NewReferencePrice(1, 0, "supplier")
	if err != referenceprice.ErrInvalidReferencePrice {
		t.Errorf("expected %v, got %v", referenceprice.ErrInvalidReferencePrice, err)
	}
}

func TestNewReferencePrice_NegativePrice_ReturnsErrInvalidReferencePrice(t *testing.T) {
	_, err := referenceprice.NewReferencePrice(1, -10, "supplier")
	if err != referenceprice.ErrInvalidReferencePrice {
		t.Errorf("expected %v, got %v", referenceprice.ErrInvalidReferencePrice, err)
	}
}
