package server

import (
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
	authHandler "pvzService/internal/services/auth/delivery/http"
)

type RouterParams struct {
	fx.In

	Logger      *slog.Logger
	AuthHandler *authHandler.Handler
}

type Router struct {
	handler *mux.Router
}

func NewRouter(p RouterParams) *Router {
	api := mux.NewRouter().PathPrefix("/api").Subrouter()

	v1 := api.PathPrefix("/v1").Subrouter()

	v1.HandleFunc("/dummyLogin", p.AuthHandler.DummyLogin).Methods(http.MethodPost, http.MethodOptions)
	v1.HandleFunc("/register", p.AuthHandler.Register).Methods(http.MethodPost, http.MethodOptions)
	v1.HandleFunc("/login", p.AuthHandler.Login).Methods(http.MethodPost, http.MethodOptions)

	router := &Router{
		handler: api,
	}

	p.Logger.Info("registered router")

	return router
}
