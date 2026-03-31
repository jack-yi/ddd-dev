package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/yangboyi/ddd-dev/user-center/internal/application"
)

type contextKey string

const (
	CtxUserID contextKey = "userId"
	CtxEmail  contextKey = "email"
	CtxRoles  contextKey = "roles"
)

func AuthMiddleware(jwtSecret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, `{"code":401,"message":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			claims, err := application.ParseToken(jwtSecret, tokenStr)
			if err != nil {
				http.Error(w, `{"code":401,"message":"invalid token"}`, http.StatusUnauthorized)
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, CtxUserID, claims.UserID)
			ctx = context.WithValue(ctx, CtxEmail, claims.Email)
			ctx = context.WithValue(ctx, CtxRoles, claims.Roles)
			next(w, r.WithContext(ctx))
		}
	}
}

func GetUserID(ctx context.Context) int64 {
	v, _ := ctx.Value(CtxUserID).(int64)
	return v
}

func GetRoles(ctx context.Context) []string {
	v, _ := ctx.Value(CtxRoles).([]string)
	return v
}

func HasRole(ctx context.Context, role string) bool {
	roles := GetRoles(ctx)
	for _, r := range roles {
		if r == role || r == "super_admin" {
			return true
		}
	}
	return false
}
