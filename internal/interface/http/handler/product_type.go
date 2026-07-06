package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-service/internal/interface/http/dto"
	producttypeinterface "stock-service/internal/interface/product_type"
)

type ProductTypeHandler struct {
	adapter *producttypeinterface.Adapter
}

func NewProductTypeHandler(adapter *producttypeinterface.Adapter) *ProductTypeHandler {
	return &ProductTypeHandler{adapter: adapter}
}

func (h *ProductTypeHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/products/{productId}/types", h.Create)
	mux.HandleFunc("GET /api/v1/products/{productId}/types", h.List)
}

func (h *ProductTypeHandler) Create(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.ParseInt(r.PathValue("productId"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}
	var body dto.CreateProductTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.Create(producttypeinterface.CreateTypeParams{
		ProductID: int32(productID), Name: body.Name,
		Value: body.Value, SortOrder: body.SortOrder,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusCreated, result)
}

func (h *ProductTypeHandler) List(w http.ResponseWriter, r *http.Request) {
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
