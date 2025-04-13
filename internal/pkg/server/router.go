package server

import (
	"github.com/gorilla/mux"
	"github.com/marrgancovka/pvzService/internal/pkg/middleware"
	authHandler "github.com/marrgancovka/pvzService/internal/services/auth/delivery/http"
	pvzHandler "github.com/marrgancovka/pvzService/internal/services/pvz/delivery/http"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
)

type RouterParams struct {
	fx.In

	Logger         *slog.Logger
	AuthHandler    *authHandler.Handler
	PvzHandler     *pvzHandler.Handler
	AuthMiddleware *middleware.AuthMiddleware
}

type Router struct {
	handler *mux.Router
}

func NewRouter(p RouterParams) *Router {
	api := mux.NewRouter().PathPrefix("/api").Subrouter()
	api.Use(middleware.CORSMiddleware)

	v1 := api.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/dummyLogin", p.AuthHandler.DummyLogin).Methods(http.MethodPost, http.MethodOptions)
	v1.HandleFunc("/register", p.AuthHandler.Register).Methods(http.MethodPost, http.MethodOptions)
	v1.HandleFunc("/login", p.AuthHandler.Login).Methods(http.MethodPost, http.MethodOptions)

	pvzGrpc := v1.PathPrefix("/pvzGrpc").Subrouter()
	pvzGrpc.HandleFunc("", p.PvzHandler.GetPvzList).Methods(http.MethodGet, http.MethodOptions)

	pvz := v1.PathPrefix("/pvz").Subrouter()
	pvz.Use(p.AuthMiddleware.AuthMiddleware)
	pvz.HandleFunc("", p.PvzHandler.CreatePvz).Methods(http.MethodPost, http.MethodOptions)
	pvz.HandleFunc("", p.PvzHandler.GetPvzs).Methods(http.MethodGet, http.MethodOptions)
	pvz.HandleFunc("/{pvzId}/close_last_reception", p.PvzHandler.CloseLastReception).Methods(http.MethodPost, http.MethodOptions)
	pvz.HandleFunc("/{pvzId}/delete_last_product", p.PvzHandler.DeleteLastProduct).Methods(http.MethodPost, http.MethodOptions)

	reception := v1.PathPrefix("/receptions").Subrouter()
	reception.Use(p.AuthMiddleware.AuthMiddleware)
	reception.HandleFunc("/receptions", p.PvzHandler.CreateReception).Methods(http.MethodPost, http.MethodOptions)

	product := v1.PathPrefix("/products").Subrouter()
	product.Use(p.AuthMiddleware.AuthMiddleware)
	product.HandleFunc("", p.PvzHandler.AddProduct).Methods(http.MethodPost, http.MethodOptions)

	router := &Router{
		handler: api,
	}

	p.Logger.Info("registered router")

	return router
}
