package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yangboyi/ddd-dev/backend/internal/application"
	"github.com/yangboyi/ddd-dev/backend/internal/model/dto"
	"github.com/yangboyi/ddd-dev/backend/pkg/response"
)

type PublishHandler struct {
	app *application.PublishApp
}

func NewPublishHandler(app *application.PublishApp) *PublishHandler {
	return &PublishHandler{app: app}
}

func (h *PublishHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePublishTaskReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, 400, "invalid request body")
		return
	}

	task, err := h.app.CreateTask(r.Context(), &req)
	if err != nil {
		response.Error(w, 500, err.Error())
		return
	}
	response.Success(w, task)
}

func (h *PublishHandler) List(w http.ResponseWriter, r *http.Request) {
	var filter dto.PublishTaskFilter
	filter.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	filter.PageSize, _ = strconv.Atoi(r.URL.Query().Get("pageSize"))

	if v := r.URL.Query().Get("status"); v != "" {
		filter.Status = &v
	}
	if v := r.URL.Query().Get("targetPlatform"); v != "" {
		filter.TargetPlatform = &v
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
