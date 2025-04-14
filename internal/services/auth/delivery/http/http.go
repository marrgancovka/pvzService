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
	const op = "auth.Handler.DummyLogin"
	logger := h.logger.With("op", op)

	var role *models.DummyLogin

	if err := reader.ReadRequestData(r, &role); err != nil {
		logger.Error("error read request data: " + err.Error())
		responser.SendErr(w, http.StatusBadRequest, auth.ErrBadRequest.Error())
		return
	}

	if role == nil {
		logger.Error("role is nil")
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
			responser.SendErr(w, http.StatusInternalServerError, "internal mainServer error")
			return
		}
	}

	logger.Info("success dummy login: " + token)
	responser.SendOk(w, http.StatusOK, token)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "auth.Handler.Login"
	logger := h.logger.With("op", op)

	userData := &models.Users{}
	if err := reader.ReadRequestData(r, userData); err != nil {
		logger.Error("error read request data: " + err.Error())
		responser.SendErr(w, http.StatusBadRequest, auth.ErrBadRequest.Error())
		return
	}

	if userData.Email == "" || userData.Password == "" {
		logger.Error("email or password is empty")
		responser.SendErr(w, http.StatusBadRequest, auth.ErrBadRequest.Error())
		return
	}

	token, err := h.usecase.Login(r.Context(), userData)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound) || errors.Is(err, auth.ErrIncorrectPasswordOrEmail):
			responser.SendErr(w, http.StatusBadRequest, auth.ErrIncorrectPasswordOrEmail.Error())
			return
		default:
			responser.SendErr(w, http.StatusInternalServerError, "internal mainServer error")
			return
		}
	}

	logger.Info("success login user: " + token)
	responser.SendOk(w, http.StatusOK, token)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "auth.Handler.Register"
	logger := h.logger.With("op", op)

	userData := &models.Users{}
	if err := reader.ReadRequestData(r, userData); err != nil {
		logger.Error("error read request data: " + err.Error())
		responser.SendErr(w, http.StatusBadRequest, auth.ErrBadRequest.Error())
		return
	}

	if userData.Email == "" || userData.Password == "" || userData.Role == "" {
		logger.Error("email, password or role is empty")
		responser.SendErr(w, http.StatusBadRequest, auth.ErrBadRequest.Error())
		return
	}

	token, err := h.usecase.Register(r.Context(), userData)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserAlreadyExists):
			responser.SendErr(w, http.StatusBadRequest, auth.ErrUserAlreadyExists.Error())
			return
		case errors.Is(err, auth.ErrIncorrectRole):
			responser.SendErr(w, http.StatusBadRequest, auth.ErrIncorrectRole.Error())
			return
		default:
			responser.SendErr(w, http.StatusInternalServerError, "internal mainServer error")
			return
		}
	}

	logger.Info("success register user: " + token)
	responser.SendOk(w, http.StatusCreated, token)
}
