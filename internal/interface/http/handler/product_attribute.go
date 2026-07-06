package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-service/internal/interface/http/dto"
	productattributeinterface "stock-service/internal/interface/product_attribute"
)

type ProductAttributeHandler struct {
	adapter *productattributeinterface.Adapter
}

func NewProductAttributeHandler(adapter *productattributeinterface.Adapter) *ProductAttributeHandler {
	return &ProductAttributeHandler{adapter: adapter}
}

func (h *ProductAttributeHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/products/{productId}/attributes", h.Create)
	mux.HandleFunc("GET /api/v1/products/{productId}/attributes", h.List)
}

func (h *ProductAttributeHandler) Create(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.ParseInt(r.PathValue("productId"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}
	var body dto.CreateProductAttributeRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.Create(productattributeinterface.CreateAttributeParams{
		ProductID: int32(productID), Key: body.Key, Value: body.Value,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusCreated, result)
}

func (h *ProductAttributeHandler) List(w http.ResponseWriter, r *http.Request) {
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
