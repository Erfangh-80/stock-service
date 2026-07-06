package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-service/internal/interface/http/dto"
	productimageinterface "stock-service/internal/interface/product_image"
)

type ProductImageHandler struct {
	adapter *productimageinterface.Adapter
}

func NewProductImageHandler(adapter *productimageinterface.Adapter) *ProductImageHandler {
	return &ProductImageHandler{adapter: adapter}
}

func (h *ProductImageHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/products/{productId}/images", h.Create)
	mux.HandleFunc("GET /api/v1/products/{productId}/images", h.List)
	mux.HandleFunc("DELETE /api/v1/products/images/{id}", h.Delete)
}

func (h *ProductImageHandler) Create(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.ParseInt(r.PathValue("productId"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}
	var body dto.CreateProductImageRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.Create(productimageinterface.CreateImageParams{
		ProductID: int32(productID), FileID: body.FileID, SortOrder: body.SortOrder,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusCreated, result)
}

func (h *ProductImageHandler) List(w http.ResponseWriter, r *http.Request) {
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

func (h *ProductImageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	if err := h.adapter.Delete(id); err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
