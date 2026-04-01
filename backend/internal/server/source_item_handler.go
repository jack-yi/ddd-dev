package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yangboyi/ddd-dev/backend/internal/application"
	"github.com/yangboyi/ddd-dev/backend/internal/model/dto"
	"github.com/yangboyi/ddd-dev/backend/pkg/response"
)

type SourceItemHandler struct {
	app *application.SourceItemApp
}

func NewSourceItemHandler(app *application.SourceItemApp) *SourceItemHandler {
	return &SourceItemHandler{app: app}
}

func (h *SourceItemHandler) Import(w http.ResponseWriter, r *http.Request) {
	var req dto.ImportSourceItemReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, 400, "invalid request body")
		return
	}

	item, err := h.app.Import(r.Context(), &req)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, item)
}

func (h *SourceItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, 400, "invalid id")
		return
	}

	item, err := h.app.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, item)
}

func (h *SourceItemHandler) List(w http.ResponseWriter, r *http.Request) {
	var filter dto.SourceItemFilter
	filter.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	filter.PageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))

	if v := r.URL.Query().Get("platform"); v != "" {
		filter.Platform = &v
	}
	if v := r.URL.Query().Get("category"); v != "" {
		filter.Category = &v
	}
	if v := r.URL.Query().Get("status"); v != "" {
		filter.Status = &v
	}
	if v := r.URL.Query().Get("keyword"); v != "" {
		filter.Keyword = &v
	}

	result, err := h.app.List(r.Context(), &filter)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, map[string]interface{}{
		"items": result.Items,
		"total": result.Total,
	})
}

func (h *SourceItemHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, 400, "invalid id")
		return
	}

	var req dto.UpdateSourceItemStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, 400, "invalid request body")
		return
	}

	if err := h.app.UpdateStatus(r.Context(), id, req.Status); err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, nil)
}

func (h *SourceItemHandler) AddTag(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(w, 400, "invalid id")
		return
	}

	var req dto.AddTagReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, 400, "invalid request body")
		return
	}

	if err := h.app.AddTag(r.Context(), id, req.Tag); err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, nil)
}
