package api

import (
	"net/http"

	"github.com/yangboyi/ddd-dev/user-center/internal/application"
	"github.com/yangboyi/ddd-dev/user-center/internal/middleware"
	"github.com/yangboyi/ddd-dev/user-center/internal/model/dto"
)

type InitHandler struct {
	userApp *application.UserApp
}

func NewInitHandler(userApp *application.UserApp) *InitHandler {
	return &InitHandler{userApp: userApp}
}

func (h *InitHandler) Check(w http.ResponseWriter, r *http.Request) {
	needInit, err := h.userApp.NeedInit(r.Context())
	if err != nil {
		writeJSON(w, 500, err.Error())
		return
	}
	writeSuccess(w, dto.CheckInitResp{NeedInit: needInit})
}

func (h *InitHandler) InitSuperAdmin(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if err := h.userApp.InitSuperAdmin(r.Context(), userID); err != nil {
		writeJSON(w, 500, err.Error())
		return
	}
	writeSuccess(w, nil)
}
