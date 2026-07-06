package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-service/internal/interface/http/dto"
	productbundleinterface "stock-service/internal/interface/product_bundle"
)

type ProductBundleHandler struct {
	adapter *productbundleinterface.Adapter
}

func NewProductBundleHandler(adapter *productbundleinterface.Adapter) *ProductBundleHandler {
	return &ProductBundleHandler{adapter: adapter}
}

func (h *ProductBundleHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/products/{productId}/bundles", h.Create)
	mux.HandleFunc("GET /api/v1/products/{productId}/bundles", h.List)
}

func (h *ProductBundleHandler) Create(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.ParseInt(r.PathValue("productId"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}
	var body dto.CreateProductBundleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.Create(productbundleinterface.CreateBundleParams{
		ProductID: int32(productID), RelatedProductID: body.RelatedProductID,
		Type: body.Type, SortOrder: body.SortOrder,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusCreated, result)
}

func (h *ProductBundleHandler) List(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.ParseInt(r.PathValue("productId"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}
	result, err := h.adapter.List(int32(productID))
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}
