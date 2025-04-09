package http

import (
	"go.uber.org/fx"
	"log/slog"
	"net/http"
)

type Params struct {
	fx.In

	Logger *slog.Logger
	//Usecase
}

type Handler struct {
	logger *slog.Logger
	//usecase
}

func NewHandler(params Params) *Handler {
	return &Handler{
		logger: params.Logger,
		//usecase: params.Usecase,
	}
}

func (h *Handler) DummyLogin(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) Login(w http.ResponseWriter, r *http.Request)      {}
func (h *Handler) Register(w http.ResponseWriter, r *http.Request)   {}
