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
	fn func(string) (*promotion.Promotion, error)
}

func (m *mockCreatePromotion) Execute(title string) (*promotion.Promotion, error) { return m.fn(title) }

type mockGetPromotion struct {
	fn func(apppromotion.GetPromotionInput) (*promotion.Promotion, error)
}

func (m *mockGetPromotion) Execute(i apppromotion.GetPromotionInput) (*promotion.Promotion, error) { return m.fn(i) }

type mockActivatePromotion struct {
	fn func(int64) error
}

func (m *mockActivatePromotion) Execute(id int64) error { return m.fn(id) }

type mockDeactivatePromotion struct {
	fn func(int64) error
}

func (m *mockDeactivatePromotion) Execute(id int64) error { return m.fn(id) }

func now() time.Time {
	t, _ := time.Parse(time.RFC3339, "2026-07-05T12:00:00Z")
	return t
}

func basePromotion() *promotion.Promotion {
	return &promotion.Promotion{
		ID: 1, Title: "Test Promotion",
		Status:    promotion.PromotionStatusInactive,
		CreatedAt: now(),
	}
}

func TestHandler_Create_Success(t *testing.T) {
	p := basePromotion()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{func(title string) (*promotion.Promotion, error) {
			if title != "Summer Sale" {
				t.Error("unexpected title")
			}
			return p, nil
		}},
		&mockGetPromotion{nil},
		&mockActivatePromotion{nil},
		&mockDeactivatePromotion{nil},
	)
	mux := http.NewServeMux()
	handler.NewPromotionHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/promotions", strings.NewReader(`{"title":"Summer Sale"}`))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var resp promotioninterface.PromotionResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.Title != "Test Promotion" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestHandler_Get_Success(t *testing.T) {
	p := basePromotion()
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil},
		&mockGetPromotion{func(i apppromotion.GetPromotionInput) (*promotion.Promotion, error) {
			if i.ID != 1 {
				t.Error("unexpected id")
			}
			return p, nil
		}},
		&mockActivatePromotion{nil},
		&mockDeactivatePromotion{nil},
	)
	mux := http.NewServeMux()
	handler.NewPromotionHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/promotions/1", nil)
	mux.ServeHTTP(w, r)

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
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil},
		&mockGetPromotion{func(i apppromotion.GetPromotionInput) (*promotion.Promotion, error) {
			return nil, promotion.ErrPromotionNotFound
		}},
		&mockActivatePromotion{nil},
		&mockDeactivatePromotion{nil},
	)
	mux := http.NewServeMux()
	handler.NewPromotionHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/promotions/999", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestHandler_Activate_Success(t *testing.T) {
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil},
		&mockGetPromotion{nil},
		&mockActivatePromotion{func(id int64) error {
			if id != 1 {
				t.Error("unexpected id")
			}
			return nil
		}},
		&mockDeactivatePromotion{nil},
	)
	mux := http.NewServeMux()
	handler.NewPromotionHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/promotions/1/activate", nil)
	mux.ServeHTTP(w, r)

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
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil},
		&mockGetPromotion{nil},
		&mockActivatePromotion{nil},
		&mockDeactivatePromotion{func(id int64) error {
			if id != 1 {
				t.Error("unexpected id")
			}
			return nil
		}},
	)
	mux := http.NewServeMux()
	handler.NewPromotionHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/promotions/1/deactivate", nil)
	mux.ServeHTTP(w, r)

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

func TestHandler_InvalidJSON(t *testing.T) {
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil},
		&mockGetPromotion{nil},
		&mockActivatePromotion{nil},
		&mockDeactivatePromotion{nil},
	)
	mux := http.NewServeMux()
	handler.NewPromotionHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/promotions", strings.NewReader(`{invalid}`))
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	b, _ := io.ReadAll(w.Body)
	var errResp struct{ Error string `json:"error"` }
	json.Unmarshal(b, &errResp)
	if errResp.Error != "invalid JSON" {
		t.Errorf("expected 'invalid JSON', got %q", errResp.Error)
	}
}

func TestHandler_InvalidID(t *testing.T) {
	adapter := promotioninterface.NewAdapter(
		&mockCreatePromotion{nil},
		&mockGetPromotion{nil},
		&mockActivatePromotion{nil},
		&mockDeactivatePromotion{nil},
	)
	mux := http.NewServeMux()
	handler.NewPromotionHandler(adapter).Register(mux)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/promotions/abc", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	b, _ := io.ReadAll(w.Body)
	var errResp struct{ Error string `json:"error"` }
	json.Unmarshal(b, &errResp)
	if errResp.Error != "invalid id" {
		t.Errorf("expected 'invalid id', got %q", errResp.Error)
	}
}
