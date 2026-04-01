package api

import (
	"net/http"

	"github.com/yangboyi/ddd-dev/user-center/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
)

func RegisterRoutes(engine *rest.Server, auth *AuthHandler, user *UserHandler, init *InitHandler, jwtSecret string) {
	authMw := middleware.AuthMiddleware(jwtSecret)

	// Public routes
	engine.AddRoutes([]rest.Route{
		{Method: http.MethodGet, Path: "/api/auth/google/login", Handler: auth.GoogleLogin},
		{Method: http.MethodGet, Path: "/api/auth/google/callback", Handler: auth.GoogleCallback},
		{Method: http.MethodPost, Path: "/api/auth/login", Handler: auth.PasswordLogin},
		{Method: http.MethodGet, Path: "/api/init/check", Handler: init.Check},
	})

	// Protected routes
	engine.AddRoutes(rest.WithMiddleware(authMw,
		rest.Route{Method: http.MethodGet, Path: "/api/auth/me", Handler: auth.Me},
		rest.Route{Method: http.MethodPost, Path: "/api/init/super-admin", Handler: init.InitSuperAdmin},
		rest.Route{Method: http.MethodGet, Path: "/api/users", Handler: user.List},
		rest.Route{Method: http.MethodPut, Path: "/api/users/status", Handler: user.UpdateStatus},
		rest.Route{Method: http.MethodPut, Path: "/api/users/role", Handler: user.AssignRole},
	))
}
