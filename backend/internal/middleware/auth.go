package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/yangboyi/ddd-dev/user-center/proto/pb"
	"github.com/zeromicro/go-zero/zrpc"
)

type contextKey string

const (
	CtxUserID contextKey = "userId"
	CtxRoles  contextKey = "roles"
)

type AuthMiddleware struct {
	userCenterClient pb.UserCenterClient
}

func NewAuthMiddleware(conn zrpc.Client) *AuthMiddleware {
	return &AuthMiddleware{
		userCenterClient: pb.NewUserCenterClient(conn.Conn()),
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, `{"code":401,"message":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		resp, err := m.userCenterClient.VerifyToken(r.Context(), &pb.VerifyTokenReq{Token: tokenStr})
		if err != nil {
			http.Error(w, `{"code":401,"message":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, CtxUserID, resp.UserId)
		ctx = context.WithValue(ctx, CtxRoles, resp.Roles)
		next(w, r.WithContext(ctx))
	}
}
