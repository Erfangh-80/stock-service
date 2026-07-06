package promotioninterface_test

import (
	"testing"
	"time"

	apppromotion "stock-service/internal/application/promotion"
	"stock-service/internal/domain/promotion"
	iface "stock-service/internal/interface"
	promotioninterface "stock-service/internal/interface/promotion"
)

type mockCreatePromotion struct{ fn func(apppromotion.CreatePromotionInput) (*promotion.Promotion, error) }

func (m *mockCreatePromotion) Execute(i apppromotion.CreatePromotionInput) (*promotion.Promotion, error) {
	return m.fn(i)
}

type mockGetPromotion struct{ fn func(apppromotion.GetPromotionInput) (*promotion.Promotion, error) }

func (m *mockGetPromotion) Execute(i apppromotion.GetPromotionInput) (*promotion.Promotion, error) {
	return m.fn(i)
}

type mockActivatePromotion struct{ fn func(int64) error }

func (m *mockActivatePromotion) Execute(id int64) error { return m.fn(id) }

type mockDeactivatePromotion struct{ fn func(int64) error }

func (m *mockDeactivatePromotion) Execute(id int64) error { return m.fn(id) }

type mockUpdatePromotion struct{ fn func(apppromotion.UpdatePromotionInput) (*promotion.Promotion, error) }

func (m *mockUpdatePromotion) Execute(i apppromotion.UpdatePromotionInput) (*promotion.Promotion, error) {
	return m.fn(i)
}

type mockDeletePromotion struct{ fn func(apppromotion.DeletePromotionInput) error }

func (m *mockDeletePromotion) Execute(i apppromotion.DeletePromotionInput) error { return m.fn(i) }

type mockListPromotions struct{ fn func(apppromotion.ListPromotionsInput) (*apppromotion.ListPromotionsOutput, error) }

func (m *mockListPromotions) Execute(i apppromotion.ListPromotionsInput) (*apppromotion.ListPromotionsOutput, error) {
	return m.fn(i)
}

func basePromotion() *promotion.Promotion {
	return &promotion.Promotion{
		ID: 1, Title: "Summer Sale",
		DiscountType:  promotion.DiscountTypePercentage,
		DiscountValue: 10,
		Status:        promotion.PromotionStatusInactive,
		CreatedAt:     time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC),
	}
}

func nilMocks() (*mockUpdatePromotion, *mockDeletePromotion, *mockListPromotions) {
	return &mockUpdatePromotion{nil}, &mockDeletePromotion{nil}, &mockListPromotions{nil}
}

func TestPromotionAdapter_Create_Success(t *testing.T) {
	p := basePromotion()
	up, del, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{func(i apppromotion.CreatePromotionInput) (*promotion.Promotion, error) { return p, nil }},
		&mockGetPromotion{nil}, &mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
		up, del, lst,
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
	up, del, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{func(i apppromotion.CreatePromotionInput) (*promotion.Promotion, error) {
			return nil, promotion.ErrTitleRequired
		}},
		&mockGetPromotion{nil}, &mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
		up, del, lst,
	)
	_, err := adapter.Create(promotioninterface.CreatePromotionParams{Title: ""})
	if err != iface.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestPromotionAdapter_Get_Success(t *testing.T) {
	p := basePromotion()
	up, del, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil},
		&mockGetPromotion{func(i apppromotion.GetPromotionInput) (*promotion.Promotion, error) { return p, nil }},
		&mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
		up, del, lst,
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
	up, del, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil},
		&mockGetPromotion{func(i apppromotion.GetPromotionInput) (*promotion.Promotion, error) {
			return nil, promotion.ErrPromotionNotFound
		}},
		&mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
		up, del, lst,
	)
	_, err := adapter.Get(999)
	if err != iface.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestPromotionAdapter_Activate_Success(t *testing.T) {
	up, del, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil}, &mockGetPromotion{nil},
		&mockActivatePromotion{func(id int64) error { return nil }},
		&mockDeactivatePromotion{nil},
		up, del, lst,
	)
	err := adapter.Activate(1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPromotionAdapter_Deactivate_Success(t *testing.T) {
	up, del, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil}, &mockGetPromotion{nil},
		&mockActivatePromotion{nil},
		&mockDeactivatePromotion{func(id int64) error { return nil }},
		up, del, lst,
	)
	err := adapter.Deactivate(1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPromotionAdapter_Update_Success(t *testing.T) {
	p := basePromotion()
	_, del, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil}, &mockGetPromotion{nil},
		&mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
		&mockUpdatePromotion{func(i apppromotion.UpdatePromotionInput) (*promotion.Promotion, error) {
			p.Title = "Updated"
			return p, nil
		}},
		del, lst,
	)
	resp, err := adapter.Update(1, promotioninterface.UpdatePromotionParams{Title: strPtr("Updated")})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Title != "Updated" {
		t.Error("expected updated title")
	}
}

func TestPromotionAdapter_Delete_Success(t *testing.T) {
	up, _, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil}, &mockGetPromotion{nil},
		&mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
		up,
		&mockDeletePromotion{func(i apppromotion.DeletePromotionInput) error { return nil }},
		lst,
	)
	err := adapter.Delete(1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPromotionAdapter_List_Success(t *testing.T) {
	p := basePromotion()
	up, del, _ := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil}, &mockGetPromotion{nil},
		&mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
		up, del,
		&mockListPromotions{func(i apppromotion.ListPromotionsInput) (*apppromotion.ListPromotionsOutput, error) {
			return &apppromotion.ListPromotionsOutput{
				Promotions: []*promotion.Promotion{p},
				Total: 1, Page: 1, Limit: 10,
			}, nil
		}},
	)
	resp, err := adapter.List(promotioninterface.ListPromotionsParams{Page: 1, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Promotions) != 1 || resp.Total != 1 {
		t.Error("expected 1 promotion")
	}
}

type unknownErr struct{}

func (e unknownErr) Error() string { return "unknown" }

func TestPromotionAdapter_UnknownError_ReturnsInternal(t *testing.T) {
	up, del, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil},
		&mockGetPromotion{func(i apppromotion.GetPromotionInput) (*promotion.Promotion, error) {
			return nil, unknownErr{}
		}},
		&mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
		up, del, lst,
	)
	_, err := adapter.Get(1)
	if err != iface.ErrInternal {
		t.Errorf("expected ErrInternal, got %v", err)
	}
}

func strPtr(s string) *string { return &s }
