package http

import (
	"go.uber.org/fx"
	"log/slog"
	"net/http"
	"pvzService/internal/models"
	"pvzService/internal/services/auth"
	"pvzService/pkg/reader"
	"pvzService/pkg/responser"
)

type Params struct {
	fx.In

	Logger  *slog.Logger
	Usecase auth.Usecase
}

type Handler struct {
	logger  *slog.Logger
	usecase auth.Usecase
}

func NewHandler(params Params) *Handler {
	return &Handler{
		logger:  params.Logger,
		usecase: params.Usecase,
	}
}

func (h *Handler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var role *models.DummyLogin

	if err := reader.ReadRequestData(r, &role); err != nil {
		responser.SendErr(w, http.StatusBadRequest, "ошибка в чтении данных")
		return
	}

	token, err := h.usecase.DummyLogin(r.Context(), role)
	if err != nil {
		// TODO: обработка ошибок
		responser.SendErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	responser.SendOk(w, http.StatusOK, token)

}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	userData := &models.Users{}
	if err := reader.ReadRequestData(r, userData); err != nil {
		responser.SendErr(w, http.StatusBadRequest, "ошибка в чтении данных")
		return
	}

	token, err := h.usecase.Login(r.Context(), userData)
	if err != nil {
		// TODO: обработка ошибок
		responser.SendErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	responser.SendOk(w, http.StatusOK, token)
}
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	userData := &models.Users{}
	if err := reader.ReadRequestData(r, userData); err != nil {
		responser.SendErr(w, http.StatusBadRequest, "ошибка в чтении данных")
		return
	}

	token, err := h.usecase.Register(r.Context(), userData)
	if err != nil {
		// TODO: обработка ошибок
		responser.SendErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	responser.SendOk(w, http.StatusCreated, token)
}
