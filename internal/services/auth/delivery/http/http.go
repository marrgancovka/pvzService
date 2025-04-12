package http

import (
	"errors"
	"github.com/marrgancovka/pvzService/internal/models"
	"github.com/marrgancovka/pvzService/internal/services/auth"
	"github.com/marrgancovka/pvzService/pkg/reader"
	"github.com/marrgancovka/pvzService/pkg/responser"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
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
		responser.SendErr(w, http.StatusBadRequest, auth.ErrBadRequest.Error())
		return
	}
	if role == nil {
		responser.SendErr(w, http.StatusBadRequest, auth.ErrBadRequest.Error())
		return
	}

	token, err := h.usecase.DummyLogin(r.Context(), role)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrIncorrectRole):
			responser.SendErr(w, http.StatusBadRequest, auth.ErrIncorrectRole.Error())
			return
		default:
			responser.SendErr(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	responser.SendOk(w, http.StatusOK, token)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	userData := &models.Users{}
	if err := reader.ReadRequestData(r, userData); err != nil {
		responser.SendErr(w, http.StatusBadRequest, auth.ErrBadRequest.Error())
		return
	}
	if userData.Email == "" || userData.Password == "" {
		responser.SendErr(w, http.StatusBadRequest, auth.ErrBadRequest.Error())
		return
	}

	token, err := h.usecase.Login(r.Context(), userData)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound) || errors.Is(err, auth.ErrIncorrectData):
			responser.SendErr(w, http.StatusBadRequest, auth.ErrIncorrectData.Error())
			return
		default:
			responser.SendErr(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}
	responser.SendOk(w, http.StatusOK, token)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	userData := &models.Users{}
	if err := reader.ReadRequestData(r, userData); err != nil {
		responser.SendErr(w, http.StatusBadRequest, "ошибка в чтении данных")
		return
	}
	if userData.Email == "" || userData.Password == "" || userData.Role == "" {
		responser.SendErr(w, http.StatusBadRequest, auth.ErrBadRequest.Error())
		return
	}

	token, err := h.usecase.Register(r.Context(), userData)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrAlreadyExists):
			responser.SendErr(w, http.StatusBadRequest, auth.ErrAlreadyExists.Error())
			return
		case errors.Is(err, auth.ErrIncorrectRole):
			responser.SendErr(w, http.StatusBadRequest, auth.ErrIncorrectRole.Error())
			return
		default:
			responser.SendErr(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}
	responser.SendOk(w, http.StatusCreated, token)
}
