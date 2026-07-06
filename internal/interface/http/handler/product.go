package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock-service/internal/interface/http/dto"
	productinterface "stock-service/internal/interface/product"
)

type ProductHandler struct {
	adapter *productinterface.Adapter
}

func NewProductHandler(adapter *productinterface.Adapter) *ProductHandler {
	return &ProductHandler{adapter: adapter}
}

func (h *ProductHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/products", h.Create)
	mux.HandleFunc("GET /api/v1/products", h.List)
	mux.HandleFunc("GET /api/v1/products/{id}", h.Get)
	mux.HandleFunc("PUT /api/v1/products/{id}", h.Update)
	mux.HandleFunc("POST /api/v1/products/{id}/activate", h.Activate)
	mux.HandleFunc("POST /api/v1/products/{id}/reject", h.Reject)
	mux.HandleFunc("POST /api/v1/products/{id}/enable", h.Enable)
	mux.HandleFunc("POST /api/v1/products/{id}/disable", h.Disable)
	mux.HandleFunc("PUT /api/v1/products/{id}/seo", h.UpdateSEO)
	mux.HandleFunc("GET /api/v1/products/my", h.MyProducts)
	mux.HandleFunc("DELETE /api/v1/products/{id}", h.SoftDelete)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.Create(productinterface.CreateProductParams{
		TitleFa:          body.TitleFa,
		TitleEn:          strPtrOrNil(body.TitleEn),
		Slug:             body.Slug,
		Description:      strPtrOrNil(body.Description),
		BrandID:          body.BrandID,
		CategoryID:       body.CategoryID,
		OwnerType:        body.OwnerType,
		OwnerID:          body.OwnerID,
		IsOriginal:       body.IsOriginal,
		MetaTitle:        strPtrOrNil(body.MetaTitle),
		MetaDescription:  strPtrOrNil(body.MetaDescription),
		IndexImageFileID: body.IndexImageFileID,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusCreated, result)
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	result, err := h.adapter.Get(int32(id))
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	var body dto.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.Update(productinterface.UpdateProductParams{
		ID:               int32(id),
		TitleFa:          body.TitleFa,
		TitleEn:          body.TitleEn,
		Slug:             body.Slug,
		Description:      body.Description,
		BrandID:          body.BrandID,
		CategoryID:       body.CategoryID,
		MetaTitle:        body.MetaTitle,
		MetaDescription:  body.MetaDescription,
		IndexImageFileID: body.IndexImageFileID,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *ProductHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	result, err := h.adapter.Activate(int32(id))
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *ProductHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	result, err := h.adapter.Reject(int32(id))
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *ProductHandler) Enable(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	result, err := h.adapter.Enable(int32(id))
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *ProductHandler) Disable(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	result, err := h.adapter.Disable(int32(id))
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *ProductHandler) UpdateSEO(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	var body dto.UpdateSEORequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid JSON"})
		return
	}
	result, err := h.adapter.UpdateSEO(int32(id), body.MetaTitle, body.MetaDescription)
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var ownerType *string
	if v := q.Get("owner_type"); v != "" {
		ownerType = &v
	}

	var ownerID *int64
	if v := q.Get("owner_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid owner_id"})
			return
		}
		ownerID = &id
	}

	var status *string
	if v := q.Get("status"); v != "" {
		status = &v
	}

	var categoryID *int64
	if v := q.Get("category_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid category_id"})
			return
		}
		categoryID = &id
	}

	var brandID *int64
	if v := q.Get("brand_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid brand_id"})
			return
		}
		brandID = &id
	}

	var search *string
	if v := q.Get("search"); v != "" {
		search = &v
	}

	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	result, err := h.adapter.List(productinterface.ListProductFilter{
		OwnerType:  ownerType,
		OwnerID:    ownerID,
		Status:     status,
		CategoryID: categoryID,
		BrandID:    brandID,
		Search:     search,
		Page:       page,
		Limit:      limit,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *ProductHandler) MyProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	ownerIDStr := q.Get("owner_id")
	if ownerIDStr == "" {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "owner_id is required"})
		return
	}

	ownerID, err := strconv.ParseInt(ownerIDStr, 10, 64)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid owner_id"})
		return
	}

	ownerType := "user"
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	result, err := h.adapter.List(productinterface.ListProductFilter{
		OwnerType: &ownerType,
		OwnerID:   &ownerID,
		Page:      page,
		Limit:     limit,
	})
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func (h *ProductHandler) SoftDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		dto.EncodeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	result, err := h.adapter.SoftDelete(int32(id))
	if err != nil {
		dto.HandleError(w, err)
		return
	}
	dto.EncodeJSON(w, http.StatusOK, result)
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
