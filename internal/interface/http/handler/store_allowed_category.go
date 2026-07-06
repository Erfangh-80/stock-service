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
	mux.HandleFunc("GET /api/v1/store-categories", h.List)
	mux.HandleFunc("GET /api/v1/store-categories/{id}", h.Get)
	mux.HandleFunc("DELETE /api/v1/store-categories/{id}", h.Delete)
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

func (h *StoreAllowedCategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	var storeID *int64
	if sid := r.URL.Query().Get("store_id"); sid != "" {
		if v, err := strconv.ParseInt(sid, 10, 64); err == nil {
			storeID = &v
		}
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20
	}

	result, err := h.adapter.List(storeallowedcategoryinterface.ListCategoriesParams{
		StoreID: storeID,
		Page:    page,
		Limit:   limit,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *StoreAllowedCategoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	result, err := h.adapter.Get(storeallowedcategoryinterface.GetCategoryParams{ID: id})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *StoreAllowedCategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	if err := h.adapter.Delete(storeallowedcategoryinterface.DeleteCategoryParams{ID: id}); err != nil {
		dto.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *StoreAllowedCategoryHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	if err := h.adapter.Approve(struct{ CategoryID int64 }{CategoryID: id}); err != nil {
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

	var body dto.RejectStoreCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	if err := h.adapter.Reject(storeallowedcategoryinterface.RejectCategoryParams{
		CategoryID:  id,
		SupportNote: body.SupportNote,
	}); err != nil {
		dto.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
