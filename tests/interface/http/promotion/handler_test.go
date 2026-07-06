package promotionhttp_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	apppromotion "stock-service/internal/application/promotion"
	"stock-service/internal/domain/promotion"
	"stock-service/internal/interface/http/handler"
	promotioninterface "stock-service/internal/interface/promotion"
)

type mockCreatePromotion struct {
	fn func(apppromotion.CreatePromotionInput) (*promotion.Promotion, error)
}

func (m *mockCreatePromotion) Execute(i apppromotion.CreatePromotionInput) (*promotion.Promotion, error) {
	return m.fn(i)
}

type mockGetPromotion struct {
	fn func(apppromotion.GetPromotionInput) (*promotion.Promotion, error)
}

func (m *mockGetPromotion) Execute(i apppromotion.GetPromotionInput) (*promotion.Promotion, error) {
	return m.fn(i)
}

type mockActivatePromotion struct {
	fn func(int64) error
}

func (m *mockActivatePromotion) Execute(id int64) error { return m.fn(id) }

type mockDeactivatePromotion struct {
	fn func(int64) error
}

func (m *mockDeactivatePromotion) Execute(id int64) error { return m.fn(id) }

type mockUpdatePromotion struct {
	fn func(apppromotion.UpdatePromotionInput) (*promotion.Promotion, error)
}

func (m *mockUpdatePromotion) Execute(i apppromotion.UpdatePromotionInput) (*promotion.Promotion, error) {
	return m.fn(i)
}

type mockDeletePromotion struct {
	fn func(apppromotion.DeletePromotionInput) error
}

func (m *mockDeletePromotion) Execute(i apppromotion.DeletePromotionInput) error { return m.fn(i) }

type mockListPromotions struct {
	fn func(apppromotion.ListPromotionsInput) (*apppromotion.ListPromotionsOutput, error)
}

func (m *mockListPromotions) Execute(i apppromotion.ListPromotionsInput) (*apppromotion.ListPromotionsOutput, error) {
	return m.fn(i)
}

func nilMocks() (*mockUpdatePromotion, *mockDeletePromotion, *mockListPromotions) {
	return &mockUpdatePromotion{nil}, &mockDeletePromotion{nil}, &mockListPromotions{nil}
}

func now() time.Time {
	t, _ := time.Parse(time.RFC3339, "2026-07-05T12:00:00Z")
	return t
}

func basePromotion() *promotion.Promotion {
	return &promotion.Promotion{
		ID: 1, Title: "Test Promotion",
		DiscountType:  promotion.DiscountTypePercentage,
		DiscountValue: 10,
		Status:        promotion.PromotionStatusInactive,
		CreatedAt:     now(),
	}
}

func doRequest(adapter *promotioninterface.Adapter, method, path, body string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	handler.NewPromotionHandler(adapter).Register(mux)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	mux.ServeHTTP(w, r)
	return w
}

func TestHandler_Create_Success(t *testing.T) {
	p := basePromotion()
	up, del, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{func(i apppromotion.CreatePromotionInput) (*promotion.Promotion, error) {
			return p, nil
		}},
		&mockGetPromotion{nil}, &mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
		up, del, lst,
	)
	w := doRequest(adapter, http.MethodPost, "/api/v1/promotions", `{"title":"Test","discount_type":"percentage","discount_value":10}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var resp promotioninterface.PromotionResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.ID != 1 || resp.Title != "Test Promotion" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestHandler_Get_Success(t *testing.T) {
	p := basePromotion()
	up, del, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil},
		&mockGetPromotion{func(i apppromotion.GetPromotionInput) (*promotion.Promotion, error) {
			return p, nil
		}},
		&mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
		up, del, lst,
	)
	w := doRequest(adapter, http.MethodGet, "/api/v1/promotions/1", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp promotioninterface.PromotionResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.ID != 1 {
		t.Errorf("expected ID 1, got %d", resp.ID)
	}
}

func TestHandler_Get_NotFound(t *testing.T) {
	up, del, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil},
		&mockGetPromotion{func(i apppromotion.GetPromotionInput) (*promotion.Promotion, error) {
			return nil, promotion.ErrPromotionNotFound
		}},
		&mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
		up, del, lst,
	)
	w := doRequest(adapter, http.MethodGet, "/api/v1/promotions/999", "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestHandler_Activate_Success(t *testing.T) {
	up, del, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil}, &mockGetPromotion{nil},
		&mockActivatePromotion{func(id int64) error { return nil }},
		&mockDeactivatePromotion{nil},
		up, del, lst,
	)
	w := doRequest(adapter, http.MethodPost, "/api/v1/promotions/1/activate", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	b, _ := io.ReadAll(w.Body)
	var resp map[string]string
	json.Unmarshal(b, &resp)
	if resp["status"] != "activated" {
		t.Errorf("expected 'activated', got %q", resp["status"])
	}
}

func TestHandler_Deactivate_Success(t *testing.T) {
	up, del, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil}, &mockGetPromotion{nil},
		&mockActivatePromotion{nil},
		&mockDeactivatePromotion{func(id int64) error { return nil }},
		up, del, lst,
	)
	w := doRequest(adapter, http.MethodPost, "/api/v1/promotions/1/deactivate", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	b, _ := io.ReadAll(w.Body)
	var resp map[string]string
	json.Unmarshal(b, &resp)
	if resp["status"] != "deactivated" {
		t.Errorf("expected 'deactivated', got %q", resp["status"])
	}
}

func TestHandler_Update_Success(t *testing.T) {
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
	w := doRequest(adapter, http.MethodPut, "/api/v1/promotions/1", `{"title":"Updated"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHandler_Delete_Success(t *testing.T) {
	up, _, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil}, &mockGetPromotion{nil},
		&mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
		up,
		&mockDeletePromotion{func(i apppromotion.DeletePromotionInput) error { return nil }},
		lst,
	)
	w := doRequest(adapter, http.MethodDelete, "/api/v1/promotions/1", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHandler_List_Success(t *testing.T) {
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
	w := doRequest(adapter, http.MethodGet, "/api/v1/promotions?page=1&limit=10", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp promotioninterface.PromotionListResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp.Promotions) != 1 || resp.Total != 1 {
		t.Error("expected 1 promotion")
	}
}

func TestHandler_InvalidJSON(t *testing.T) {
	up, del, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil}, &mockGetPromotion{nil},
		&mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
		up, del, lst,
	)
	w := doRequest(adapter, http.MethodPost, "/api/v1/promotions", `{invalid}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandler_InvalidID(t *testing.T) {
	up, del, lst := nilMocks()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil}, &mockGetPromotion{nil},
		&mockActivatePromotion{nil}, &mockDeactivatePromotion{nil},
		up, del, lst,
	)
	w := doRequest(adapter, http.MethodGet, "/api/v1/promotions/abc", "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
