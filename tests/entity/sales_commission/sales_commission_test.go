package salescommission

import (
	"testing"

	"stock-service/internal/domain/sales_commission"
)

func TestNewSalesCommission_ValidInputs_SetsDefaults(t *testing.T) {
	sc, err := salescommission.NewSalesCommission(10, 20, salescommission.SaleModelRetail, 5.0, 100.0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sc.InventoryID != 10 {
		t.Errorf("expected InventoryID %d, got %d", 10, sc.InventoryID)
	}
	if sc.CategoryCommissionRuleID != 20 {
		t.Errorf("expected CategoryCommissionRuleID %d, got %d", 20, sc.CategoryCommissionRuleID)
	}
	if sc.SaleModel != salescommission.SaleModelRetail {
		t.Errorf("expected SaleModel %q, got %q", salescommission.SaleModelRetail, sc.SaleModel)
	}
	if sc.RatePercent != 5.0 {
		t.Errorf("expected RatePercent %f, got %f", 5.0, sc.RatePercent)
	}
	if sc.MinPrice != 100.0 {
		t.Errorf("expected MinPrice %f, got %f", 100.0, sc.MinPrice)
	}
	if sc.MinQty != nil {
		t.Errorf("expected MinQty to be nil, got %v", *sc.MinQty)
	}
	if sc.MaxPrice != nil {
		t.Errorf("expected MaxPrice to be nil, got %v", *sc.MaxPrice)
	}
}

func TestNewSalesCommission_RatePercentOver100_ReturnsErrInvalidRatePercent(t *testing.T) {
	_, err := salescommission.NewSalesCommission(1, 1, salescommission.SaleModelRetail, 150, 100)
	if err != salescommission.ErrInvalidRatePercent {
		t.Errorf("expected %v, got %v", salescommission.ErrInvalidRatePercent, err)
	}
}

func TestNewSalesCommission_RatePercentBelow0_ReturnsErrInvalidRatePercent(t *testing.T) {
	_, err := salescommission.NewSalesCommission(1, 1, salescommission.SaleModelRetail, -1, 100)
	if err != salescommission.ErrInvalidRatePercent {
		t.Errorf("expected %v, got %v", salescommission.ErrInvalidRatePercent, err)
	}
}

func TestNewSalesCommission_ZeroMinPrice_ReturnsErrInvalidMinPrice(t *testing.T) {
	_, err := salescommission.NewSalesCommission(1, 1, salescommission.SaleModelRetail, 5, 0)
	if err != salescommission.ErrInvalidMinPrice {
		t.Errorf("expected %v, got %v", salescommission.ErrInvalidMinPrice, err)
	}
}

func TestUpdateMaxPrice_ValidPrice_Succeeds(t *testing.T) {
	sc, _ := salescommission.NewSalesCommission(1, 1, salescommission.SaleModelRetail, 5, 100)
	err := sc.UpdateMaxPrice(200)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sc.MaxPrice == nil {
		t.Fatal("expected MaxPrice to be set, got nil")
	}
	if *sc.MaxPrice != 200 {
		t.Errorf("expected MaxPrice %f, got %f", 200.0, *sc.MaxPrice)
	}
}

func TestUpdateMaxPrice_PriceNotGreaterThanMin_ReturnsErrInvalidMaxPrice(t *testing.T) {
	sc, _ := salescommission.NewSalesCommission(1, 1, salescommission.SaleModelRetail, 5, 100)
	err := sc.UpdateMaxPrice(100)
	if err != salescommission.ErrInvalidMaxPrice {
		t.Errorf("expected %v, got %v", salescommission.ErrInvalidMaxPrice, err)
	}
}

func TestUpdateMinQty_ValidQty_Succeeds(t *testing.T) {
	sc, _ := salescommission.NewSalesCommission(1, 1, salescommission.SaleModelRetail, 5, 100)
	err := sc.UpdateMinQty(5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if sc.MinQty == nil {
		t.Fatal("expected MinQty to be set, got nil")
	}
	if *sc.MinQty != 5 {
		t.Errorf("expected MinQty %d, got %d", 5, *sc.MinQty)
	}
}

func TestUpdateMinQty_NegativeQty_ReturnsErrInvalidMinQty(t *testing.T) {
	sc, _ := salescommission.NewSalesCommission(1, 1, salescommission.SaleModelRetail, 5, 100)
	err := sc.UpdateMinQty(-1)
	if err != salescommission.ErrInvalidMinQty {
		t.Errorf("expected %v, got %v", salescommission.ErrInvalidMinQty, err)
	}
}
