package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yangboyi/ddd-dev/backend/internal/application"
	"github.com/yangboyi/ddd-dev/backend/internal/model/dto"
	"github.com/yangboyi/ddd-dev/backend/pkg/response"
)

type ProductHandler struct {
	app *application.ProductApp
}

func NewProductHandler(app *application.ProductApp) *ProductHandler {
	return &ProductHandler{app: app}
}

func (h *ProductHandler) CreateFromSource(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductFromSourceReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, 400, "invalid request body")
		return
	}

	product, err := h.app.CreateFromSource(r.Context(), req.SourceItemID)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, product)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, 400, "invalid id")
		return
	}

	product, err := h.app.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, product)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	var filter dto.ProductFilter
	filter.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	filter.PageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))

	if v := r.URL.Query().Get("status"); v != "" {
		filter.Status = &v
	}
	if v := r.URL.Query().Get("keyword"); v != "" {
		filter.Keyword = &v
	}

	items, total, err := h.app.List(r.Context(), &filter)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, map[string]interface{}{
		"items": items,
		"total": total,
	})
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, 400, "invalid id")
		return
	}

	var req dto.UpdateProductReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, 400, "invalid request body")
		return
	}

	product, err := h.app.Update(r.Context(), id, &req)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, product)
}

func (h *ProductHandler) MarkReady(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, 400, "invalid id")
		return
	}

	if err := h.app.MarkReady(r.Context(), id); err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, nil)
}
