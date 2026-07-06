package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-service/internal/interface/http/dto"
	salescommissioninterface "stock-service/internal/interface/sales_commission"
)

type CategoryCommissionRuleHandler struct {
	adapter *salescommissioninterface.CategoryCommissionRuleAdapter
}

func NewCategoryCommissionRuleHandler(adapter *salescommissioninterface.CategoryCommissionRuleAdapter) *CategoryCommissionRuleHandler {
	return &CategoryCommissionRuleHandler{adapter: adapter}
}

func (h *CategoryCommissionRuleHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/category-commission-rules", h.Create)
	mux.HandleFunc("GET /api/v1/category-commission-rules", h.List)
	mux.HandleFunc("GET /api/v1/category-commission-rules/{id}", h.Get)
	mux.HandleFunc("PUT /api/v1/category-commission-rules/{id}", h.Update)
	mux.HandleFunc("DELETE /api/v1/category-commission-rules/{id}", h.Delete)
}

func (h *CategoryCommissionRuleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateCategoryCommissionRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	result, err := h.adapter.Create(salescommissioninterface.CreateCategoryCommissionRuleParams{
		CategoryID:  body.CategoryID,
		RatePercent: body.RatePercent,
		MinPrice:    body.MinPrice,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusCreated, result)
}

func (h *CategoryCommissionRuleHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	result, err := h.adapter.Get(id)
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *CategoryCommissionRuleHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var categoryID *int32
	if s := q.Get("category_id"); s != "" {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid category_id"})
			return
		}
		v32 := int32(v)
		categoryID = &v32
	}

	var isActive *bool
	if s := q.Get("is_active"); s != "" {
		v := s == "true"
		isActive = &v
	}

	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	result, err := h.adapter.List(salescommissioninterface.ListCategoryCommissionRulesParams{
		CategoryID: categoryID,
		IsActive:   isActive,
		Page:       page,
		Limit:      limit,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *CategoryCommissionRuleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var body dto.UpdateCategoryCommissionRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}

	result, err := h.adapter.Update(id, salescommissioninterface.UpdateCategoryCommissionRuleParams{
		RatePercent: body.RatePercent,
		MinPrice:    body.MinPrice,
		MaxPrice:    body.MaxPrice,
		Activate:    body.Activate,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}

	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *CategoryCommissionRuleHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
