package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yangboyi/ddd-dev/user-center/internal/application"
	"github.com/yangboyi/ddd-dev/user-center/internal/config"
	"github.com/yangboyi/ddd-dev/user-center/internal/middleware"
	"github.com/yangboyi/ddd-dev/user-center/internal/model/dto"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	userApp     *application.UserApp
	oauthConfig *oauth2.Config
}

func NewAuthHandler(userApp *application.UserApp, googleCfg config.GoogleConfig) *AuthHandler {
	return &AuthHandler{
		userApp: userApp,
		oauthConfig: &oauth2.Config{
			ClientID:     googleCfg.ClientID,
			ClientSecret: googleCfg.ClientSecret,
			RedirectURL:  googleCfg.RedirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := h.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		writeJSON(w, 400, "missing code")
		return
	}

	token, err := h.oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		writeJSON(w, 500, fmt.Sprintf("exchange token: %v", err))
		return
	}

	client := h.oauthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		writeJSON(w, 500, fmt.Sprintf("get user info: %v", err))
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.Unmarshal(body, &googleUser); err != nil {
		writeJSON(w, 500, "parse google user info failed")
		return
	}

	user, jwtToken, err := h.userApp.LoginOrRegister(r.Context(), application.GoogleUserInfo{
		GoogleID: googleUser.ID,
		Email:    googleUser.Email,
		Name:     googleUser.Name,
		Avatar:   googleUser.Picture,
	})
	if err != nil {
		writeJSON(w, 500, err.Error())
		return
	}

	_ = user
	frontendURL := fmt.Sprintf("http://localhost:3000/login/callback?token=%s", jwtToken)
	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	user, err := h.userApp.GetUserInfo(r.Context(), userID)
	if err != nil {
		writeJSON(w, 500, err.Error())
		return
	}
	writeSuccess(w, dto.UserResp{
		ID: user.ID, Email: user.Email, Name: user.Name,
		Avatar: user.Avatar, Status: user.Status, Roles: user.RoleNames(),
	})
}

func (h *AuthHandler) PasswordLogin(w http.ResponseWriter, r *http.Request) {
	var req dto.PasswordLoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, 400, "invalid request body")
		return
	}

	user, token, err := h.userApp.LoginByPassword(r.Context(), req.Username, req.Password)
	if err != nil {
		writeJSON(w, 401, err.Error())
		return
	}

	writeSuccess(w, dto.LoginResp{
		Token: token,
		User: dto.UserResp{
			ID: user.ID, Email: user.Email, Name: user.Name,
			Avatar: user.Avatar, Status: user.Status, Roles: user.RoleNames(),
		},
	})
}

func writeJSON(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"code": code, "message": msg})
}

func writeSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"code": 0, "message": "ok", "data": data})
}
