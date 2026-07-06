package promotioninterface_test

import (
	"testing"
	"time"

	apppromotion "stock-service/internal/application/promotion"
	"stock-service/internal/domain/promotion"
	iface "stock-service/internal/interface"
	promotioninterface "stock-service/internal/interface/promotion"
)

type mockCreatePromotion struct{ fn func(string) (*promotion.Promotion, error) }

func (m *mockCreatePromotion) Execute(title string) (*promotion.Promotion, error) { return m.fn(title) }

type mockGetPromotion struct{ fn func(apppromotion.GetPromotionInput) (*promotion.Promotion, error) }

func (m *mockGetPromotion) Execute(i apppromotion.GetPromotionInput) (*promotion.Promotion, error) {
	return m.fn(i)
}

type mockActivatePromotion struct{ fn func(int64) error }

func (m *mockActivatePromotion) Execute(id int64) error { return m.fn(id) }

type mockDeactivatePromotion struct{ fn func(int64) error }

func (m *mockDeactivatePromotion) Execute(id int64) error { return m.fn(id) }

func basePromotion() *promotion.Promotion {
	return &promotion.Promotion{
		ID: 1, Title: "Summer Sale",
		Status:    promotion.PromotionStatusInactive,
		CreatedAt: time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
	}
}

func TestPromotionAdapter_Create_Success(t *testing.T) {
	p := basePromotion()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{func(title string) (*promotion.Promotion, error) { return p, nil }},
		&mockGetPromotion{nil}, &mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
	)
	resp, err := adapter.Create(promotioninterface.CreatePromotionParams{Title: "Summer Sale"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.Title != "Summer Sale" {
		t.Error("unexpected response")
	}
}

func TestPromotionAdapter_Create_InvalidInput(t *testing.T) {
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{func(title string) (*promotion.Promotion, error) {
			return nil, promotion.ErrTitleRequired
		}},
		&mockGetPromotion{nil}, &mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
	)
	_, err := adapter.Create(promotioninterface.CreatePromotionParams{Title: ""})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestPromotionAdapter_Get_Success(t *testing.T) {
	p := basePromotion()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil},
		&mockGetPromotion{func(i apppromotion.GetPromotionInput) (*promotion.Promotion, error) { return p, nil }},
		&mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
	)
	resp, err := adapter.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 {
		t.Error("expected ID 1")
	}
}

func TestPromotionAdapter_Get_NotFound(t *testing.T) {
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil},
		&mockGetPromotion{func(i apppromotion.GetPromotionInput) (*promotion.Promotion, error) {
			return nil, promotion.ErrPromotionNotFound
		}},
		&mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
	)
	_, err := adapter.Get(999)
	if err != iface.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestPromotionAdapter_Activate_Success(t *testing.T) {
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil}, &mockGetPromotion{nil},
		&mockActivatePromotion{func(id int64) error { return nil }},
		&mockDeactivatePromotion{nil},
	)
	err := adapter.Activate(1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPromotionAdapter_Deactivate_Success(t *testing.T) {
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil}, &mockGetPromotion{nil},
		&mockActivatePromotion{nil},
		&mockDeactivatePromotion{func(id int64) error { return nil }},
	)
	err := adapter.Deactivate(1)
	if err != nil {
		t.Fatal(err)
	}
}

type unknownErr struct{}

func (e unknownErr) Error() string { return "unknown" }

func TestPromotionAdapter_UnknownError_ReturnsInternal(t *testing.T) {
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil},
		&mockGetPromotion{func(i apppromotion.GetPromotionInput) (*promotion.Promotion, error) {
			return nil, unknownErr{}
		}},
		&mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
	)
	_, err := adapter.Get(1)
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}
