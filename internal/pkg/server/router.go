package server

import (
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"log/slog"
)

type RouterParams struct {
	fx.In

	Logger *slog.Logger
}

type Router struct {
	handler *mux.Router
}

func NewRouter(p RouterParams) *Router {
	api := mux.NewRouter().PathPrefix("/api").Subrouter()

	_ = api.PathPrefix("/v1").Subrouter()

	router := &Router{
		handler: api,
	}

	p.Logger.Info("registered router")

	return router
}
