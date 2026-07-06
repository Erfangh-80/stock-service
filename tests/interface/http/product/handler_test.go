package producthttp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	appproduct "stock-service/internal/application/product"
	"stock-service/internal/domain/product"
	"stock-service/internal/interface/http/handler"
	productinterface "stock-service/internal/interface/product"
)

type mockCreateProduct struct{ fn func(appproduct.CreateProductInput) (*product.Product, error) }
func (m *mockCreateProduct) Execute(i appproduct.CreateProductInput) (*product.Product, error) { return m.fn(i) }

type mockGetProduct struct{ fn func(appproduct.GetProductInput) (*product.Product, error) }
func (m *mockGetProduct) Execute(i appproduct.GetProductInput) (*product.Product, error) { return m.fn(i) }

type mockUpdateProduct struct{ fn func(appproduct.UpdateProductInput) (*product.Product, error) }
func (m *mockUpdateProduct) Execute(i appproduct.UpdateProductInput) (*product.Product, error) { return m.fn(i) }

type mockActivateProduct struct{ fn func(appproduct.ActivateProductInput) (*product.Product, error) }
func (m *mockActivateProduct) Execute(i appproduct.ActivateProductInput) (*product.Product, error) { return m.fn(i) }

type mockRejectProduct struct{ fn func(appproduct.RejectProductInput) (*product.Product, error) }
func (m *mockRejectProduct) Execute(i appproduct.RejectProductInput) (*product.Product, error) { return m.fn(i) }

type mockSoftDeleteProduct struct{ fn func(appproduct.SoftDeleteProductInput) (*product.Product, error) }
func (m *mockSoftDeleteProduct) Execute(i appproduct.SoftDeleteProductInput) (*product.Product, error) { return m.fn(i) }

type mockEnableProduct struct{ fn func(appproduct.EnableProductInput) (*product.Product, error) }
func (m *mockEnableProduct) Execute(i appproduct.EnableProductInput) (*product.Product, error) { return m.fn(i) }

type mockDisableProduct struct{ fn func(appproduct.DisableProductInput) (*product.Product, error) }
func (m *mockDisableProduct) Execute(i appproduct.DisableProductInput) (*product.Product, error) { return m.fn(i) }

type mockUpdateSEO struct{ fn func(appproduct.UpdateSEOInput) (*product.Product, error) }
func (m *mockUpdateSEO) Execute(i appproduct.UpdateSEOInput) (*product.Product, error) { return m.fn(i) }

type mockListProducts struct{ fn func(appproduct.ListProductsInput) (*appproduct.ListProductsOutput, error) }
func (m *mockListProducts) Execute(i appproduct.ListProductsInput) (*appproduct.ListProductsOutput, error) { return m.fn(i) }

func nilMocks() (*mockEnableProduct, *mockDisableProduct, *mockUpdateSEO, *mockListProducts) {
	return &mockEnableProduct{nil}, &mockDisableProduct{nil}, &mockUpdateSEO{nil}, &mockListProducts{nil}
}

func now() time.Time {
	t, _ := time.Parse(time.RFC3339, "2026-07-05T12:00:00Z")
	return t
}

func baseProduct() *product.Product {
	return &product.Product{
		ID: 1, TitleFa: "محصول تست", BrandID: 1, CategoryID: 1,
		Status:    product.ProductStatusPending,
		OwnerType: product.OwnerTypeSystem,
		CreatedAt: now(),
		UpdatedAt: now(),
	}
}

func doRequest(adapter *productinterface.Adapter, method, path, body string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	handler.NewProductHandler(adapter).Register(mux)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	mux.ServeHTTP(w, r)
	return w
}

func assertStatus(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if w.Code != expected {
		t.Errorf("expected status %d, got %d", expected, w.Code)
	}
}

func TestHandler_Create_Success(t *testing.T) {
	p := baseProduct()
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{func(i appproduct.CreateProductInput) (*product.Product, error) {
			if i.TitleFa != "محصول جدید" {
				t.Errorf("unexpected title: %q", i.TitleFa)
			}
			return p, nil
		}},
		&mockGetProduct{nil}, &mockUpdateProduct{nil},
		&mockActivateProduct{nil}, &mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)

	w := doRequest(adapter, http.MethodPost, "/api/v1/products", `{"title_fa":"محصول جدید","brand_id":1,"category_id":1}`)
	assertStatus(t, w, http.StatusCreated)

	var resp productinterface.ProductResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 {
		t.Error("expected id 1")
	}
}

func TestHandler_Create_InvalidJSON(t *testing.T) {
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil}, &mockGetProduct{nil}, &mockUpdateProduct{nil},
		&mockActivateProduct{nil}, &mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)
	w := doRequest(adapter, http.MethodPost, "/api/v1/products", `{invalid}`)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestHandler_Create_MissingFields(t *testing.T) {
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{func(i appproduct.CreateProductInput) (*product.Product, error) {
			return nil, product.ErrTitleFaRequired
		}},
		&mockGetProduct{nil}, &mockUpdateProduct{nil},
		&mockActivateProduct{nil}, &mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)
	w := doRequest(adapter, http.MethodPost, "/api/v1/products", `{"brand_id":1,"category_id":1}`)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestHandler_Get_Success(t *testing.T) {
	p := baseProduct()
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil},
		&mockGetProduct{func(i appproduct.GetProductInput) (*product.Product, error) { return p, nil }},
		&mockUpdateProduct{nil}, &mockActivateProduct{nil},
		&mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)
	w := doRequest(adapter, http.MethodGet, "/api/v1/products/1", "")
	assertStatus(t, w, http.StatusOK)
}

func TestHandler_Get_NotFound(t *testing.T) {
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil},
		&mockGetProduct{func(i appproduct.GetProductInput) (*product.Product, error) {
			return nil, product.ErrProductNotFound
		}},
		&mockUpdateProduct{nil}, &mockActivateProduct{nil},
		&mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)
	w := doRequest(adapter, http.MethodGet, "/api/v1/products/999", "")
	assertStatus(t, w, http.StatusNotFound)
}

func TestHandler_Get_InvalidID(t *testing.T) {
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil}, &mockGetProduct{nil}, &mockUpdateProduct{nil},
		&mockActivateProduct{nil}, &mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)
	w := doRequest(adapter, http.MethodGet, "/api/v1/products/abc", "")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestHandler_Update_Success(t *testing.T) {
	p := baseProduct()
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil}, &mockGetProduct{nil},
		&mockUpdateProduct{func(i appproduct.UpdateProductInput) (*product.Product, error) { return p, nil }},
		&mockActivateProduct{nil}, &mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)
	w := doRequest(adapter, http.MethodPut, "/api/v1/products/1", `{"title_fa":"updated"}`)
	assertStatus(t, w, http.StatusOK)
}

func TestHandler_Activate_Success(t *testing.T) {
	p := baseProduct()
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil}, &mockGetProduct{nil}, &mockUpdateProduct{nil},
		&mockActivateProduct{func(i appproduct.ActivateProductInput) (*product.Product, error) { return p, nil }},
		&mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)
	w := doRequest(adapter, http.MethodPost, "/api/v1/products/1/activate", "")
	assertStatus(t, w, http.StatusOK)
}

func TestHandler_Reject_Success(t *testing.T) {
	p := baseProduct()
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil}, &mockGetProduct{nil}, &mockUpdateProduct{nil},
		&mockActivateProduct{nil},
		&mockRejectProduct{func(i appproduct.RejectProductInput) (*product.Product, error) { return p, nil }},
		&mockSoftDeleteProduct{nil},
		en, dis, seo, lst,
	)
	w := doRequest(adapter, http.MethodPost, "/api/v1/products/1/reject", "")
	assertStatus(t, w, http.StatusOK)
}

func TestHandler_SoftDelete_Success(t *testing.T) {
	p := baseProduct()
	en, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil}, &mockGetProduct{nil}, &mockUpdateProduct{nil},
		&mockActivateProduct{nil}, &mockRejectProduct{nil},
		&mockSoftDeleteProduct{func(i appproduct.SoftDeleteProductInput) (*product.Product, error) { return p, nil }},
		en, dis, seo, lst,
	)
	w := doRequest(adapter, http.MethodDelete, "/api/v1/products/1", "")
	assertStatus(t, w, http.StatusOK)
}

func TestHandler_Enable_Success(t *testing.T) {
	p := baseProduct()
	_, dis, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil}, &mockGetProduct{nil}, &mockUpdateProduct{nil},
		&mockActivateProduct{nil}, &mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		&mockEnableProduct{func(i appproduct.EnableProductInput) (*product.Product, error) { return p, nil }},
		dis, seo, lst,
	)
	w := doRequest(adapter, http.MethodPost, "/api/v1/products/1/enable", "")
	assertStatus(t, w, http.StatusOK)
}

func TestHandler_Disable_Success(t *testing.T) {
	p := baseProduct()
	en, _, seo, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil}, &mockGetProduct{nil}, &mockUpdateProduct{nil},
		&mockActivateProduct{nil}, &mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en,
		&mockDisableProduct{func(i appproduct.DisableProductInput) (*product.Product, error) { return p, nil }},
		seo, lst,
	)
	w := doRequest(adapter, http.MethodPost, "/api/v1/products/1/disable", "")
	assertStatus(t, w, http.StatusOK)
}

func TestHandler_UpdateSEO_Success(t *testing.T) {
	p := baseProduct()
	en, dis, _, lst := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil}, &mockGetProduct{nil}, &mockUpdateProduct{nil},
		&mockActivateProduct{nil}, &mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis,
		&mockUpdateSEO{func(i appproduct.UpdateSEOInput) (*product.Product, error) { return p, nil }},
		lst,
	)
	w := doRequest(adapter, http.MethodPut, "/api/v1/products/1/seo", `{"meta_title":"test title"}`)
	assertStatus(t, w, http.StatusOK)
}

func TestHandler_List_Success(t *testing.T) {
	p := baseProduct()
	en, dis, seo, _ := nilMocks()
	adapter := productinterface.NewAdapter(
		&mockCreateProduct{nil}, &mockGetProduct{nil}, &mockUpdateProduct{nil},
		&mockActivateProduct{nil}, &mockRejectProduct{nil}, &mockSoftDeleteProduct{nil},
		en, dis, seo,
		&mockListProducts{func(i appproduct.ListProductsInput) (*appproduct.ListProductsOutput, error) {
			return &appproduct.ListProductsOutput{
				Products: []*product.Product{p},
				Total:    1, Page: 1, Limit: 10,
			}, nil
		}},
	)
	w := doRequest(adapter, http.MethodGet, "/api/v1/products", "")
	assertStatus(t, w, http.StatusOK)

	var resp productinterface.ProductListResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Products) != 1 || resp.Total != 1 {
		t.Error("expected 1 product in list")
	}
}
