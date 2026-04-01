package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yangboyi/ddd-dev/user-center/internal/application"
	"github.com/yangboyi/ddd-dev/user-center/internal/middleware"
	"github.com/yangboyi/ddd-dev/user-center/internal/model/dto"
)

type UserHandler struct {
	userApp *application.UserApp
}

func NewUserHandler(userApp *application.UserApp) *UserHandler {
	return &UserHandler{userApp: userApp}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	if !middleware.HasRole(r.Context(), "admin") {
		writeJSON(w, 403, "forbidden")
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	users, total, err := h.userApp.ListUsers(r.Context(), page, pageSize)
	if err != nil {
		writeJSON(w, 500, err.Error())
		return
	}
	writeSuccess(w, map[string]interface{}{"items": users, "total": total})
}

func (h *UserHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	if !middleware.HasRole(r.Context(), "admin") {
		writeJSON(w, 403, "forbidden")
		return
	}
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	var req dto.UpdateUserStatusReq
	json.NewDecoder(r.Body).Decode(&req)
	if err := h.userApp.UpdateStatus(r.Context(), id, req.Status); err != nil {
		writeJSON(w, 500, err.Error())
		return
	}
	writeSuccess(w, nil)
}

func (h *UserHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	if !middleware.HasRole(r.Context(), "super_admin") {
		writeJSON(w, 403, "only super admin can assign roles")
		return
	}
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	var req dto.AssignRoleReq
	json.NewDecoder(r.Body).Decode(&req)
	if err := h.userApp.AssignRole(r.Context(), id, req.RoleName); err != nil {
		writeJSON(w, 500, err.Error())
		return
	}
	writeSuccess(w, nil)
}
