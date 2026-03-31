package server

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterRoutes(engine *rest.Server, si *SourceItemHandler, p *ProductHandler, pub *PublishHandler) {
	engine.AddRoutes([]rest.Route{
		{Method: http.MethodPost, Path: "/api/source-items/import", Handler: si.Import},
		{Method: http.MethodGet, Path: "/api/source-items", Handler: si.List},
		{Method: http.MethodGet, Path: "/api/source-items/detail", Handler: si.GetByID},
		{Method: http.MethodPut, Path: "/api/source-items/status", Handler: si.UpdateStatus},
		{Method: http.MethodPost, Path: "/api/source-items/tag", Handler: si.AddTag},

		{Method: http.MethodPost, Path: "/api/products/create-from-source", Handler: p.CreateFromSource},
		{Method: http.MethodGet, Path: "/api/products", Handler: p.List},
		{Method: http.MethodGet, Path: "/api/products/detail", Handler: p.GetByID},
		{Method: http.MethodPut, Path: "/api/products", Handler: p.Update},
		{Method: http.MethodPut, Path: "/api/products/ready", Handler: p.MarkReady},

		{Method: http.MethodPost, Path: "/api/publish-tasks", Handler: pub.CreateTask},
		{Method: http.MethodGet, Path: "/api/publish-tasks", Handler: pub.List},
	})
}
