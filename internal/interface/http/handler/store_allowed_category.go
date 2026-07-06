package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-service/internal/interface/http/dto"
	storeallowedcategoryinterface "stock-service/internal/interface/store_allowed_category"
)

type StoreAllowedCategoryHandler struct {
	adapter *storeallowedcategoryinterface.Adapter
}

func NewStoreAllowedCategoryHandler(adapter *storeallowedcategoryinterface.Adapter) *StoreAllowedCategoryHandler {
	return &StoreAllowedCategoryHandler{adapter: adapter}
}

func (h *StoreAllowedCategoryHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/store-categories", h.Create)
	mux.HandleFunc("POST /api/v1/store-categories/{id}/approve", h.Approve)
	mux.HandleFunc("POST /api/v1/store-categories/{id}/reject", h.Reject)
}

func (h *StoreAllowedCategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateStoreCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	result, err := h.adapter.Create(storeallowedcategoryinterface.CreateCategoryParams{
		StoreID: body.StoreID, CategoryID: body.CategoryID,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusCreated, result)
}

func (h *StoreAllowedCategoryHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	if err := h.adapter.Approve(storeallowedcategoryinterface.ApproveCategoryParams{CategoryID: id}); err != nil {
		dto.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *StoreAllowedCategoryHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	if err := h.adapter.Reject(storeallowedcategoryinterface.RejectCategoryParams{CategoryID: id}); err != nil {
		dto.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
