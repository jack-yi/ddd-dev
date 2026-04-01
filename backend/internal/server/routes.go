package server

import (
	"net/http"

	"github.com/yangboyi/ddd-dev/backend/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
)

func RegisterRoutes(engine *rest.Server, si *SourceItemHandler, p *ProductHandler, pub *PublishHandler, authMw *middleware.AuthMiddleware) {
	engine.AddRoutes(rest.WithMiddleware(authMw.Handle,
		rest.Route{Method: http.MethodPost, Path: "/api/source-items/import", Handler: si.Import},
		rest.Route{Method: http.MethodGet, Path: "/api/source-items", Handler: si.List},
		rest.Route{Method: http.MethodGet, Path: "/api/source-items/detail", Handler: si.GetByID},
		rest.Route{Method: http.MethodPut, Path: "/api/source-items/status", Handler: si.UpdateStatus},
		rest.Route{Method: http.MethodPost, Path: "/api/source-items/tag", Handler: si.AddTag},

		rest.Route{Method: http.MethodPost, Path: "/api/products/create-from-source", Handler: p.CreateFromSource},
		rest.Route{Method: http.MethodGet, Path: "/api/products", Handler: p.List},
		rest.Route{Method: http.MethodGet, Path: "/api/products/detail", Handler: p.GetByID},
		rest.Route{Method: http.MethodPut, Path: "/api/products", Handler: p.Update},
		rest.Route{Method: http.MethodPut, Path: "/api/products/ready", Handler: p.MarkReady},

		rest.Route{Method: http.MethodPost, Path: "/api/publish-tasks", Handler: pub.CreateTask},
		rest.Route{Method: http.MethodGet, Path: "/api/publish-tasks", Handler: pub.List},
	))
}
