package server

import (
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
	"pvzService/internal/pkg/middleware"
	authHandler "pvzService/internal/services/auth/delivery/http"
	pvzHandler "pvzService/internal/services/pvz/delivery/http"
)

type RouterParams struct {
	fx.In

	Logger      *slog.Logger
	AuthHandler *authHandler.Handler
	PvzHandler  *pvzHandler.Handler
}

type Router struct {
	handler *mux.Router
}

func NewRouter(p RouterParams) *Router {
	api := mux.NewRouter().PathPrefix("/api").Subrouter()

	v1 := api.PathPrefix("/v1").Subrouter()
	v1.Use(middleware.CORSMiddleware)

	v1.HandleFunc("/dummyLogin", p.AuthHandler.DummyLogin).Methods(http.MethodPost, http.MethodOptions)
	v1.HandleFunc("/register", p.AuthHandler.Register).Methods(http.MethodPost, http.MethodOptions)
	v1.HandleFunc("/login", p.AuthHandler.Login).Methods(http.MethodPost, http.MethodOptions)

	v1.HandleFunc("/pvz", p.PvzHandler.CreatePvz).Methods(http.MethodPost, http.MethodOptions)
	v1.HandleFunc("/pvz", p.PvzHandler.GetPvzList).Methods(http.MethodGet, http.MethodOptions)
	v1.HandleFunc("/pvz/{pvzId}/close_last_reception", p.PvzHandler.CloseLastReception).Methods(http.MethodPost, http.MethodOptions)
	v1.HandleFunc("/pvz/{pvzId}/delete_last_product", p.PvzHandler.DeleteLastProduct).Methods(http.MethodPost, http.MethodOptions)

	v1.HandleFunc("/receptions", p.PvzHandler.CreateReception).Methods(http.MethodPost, http.MethodOptions)
	v1.HandleFunc("/products", p.PvzHandler.AddProduct).Methods(http.MethodPost, http.MethodOptions)

	router := &Router{
		handler: api,
	}

	p.Logger.Info("registered router")

	return router
}
