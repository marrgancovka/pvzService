package http

import (
	"go.uber.org/fx"
	"log/slog"
	"net/http"
	"pvzService/internal/models"
	"pvzService/internal/services/pvz"
	"pvzService/pkg/reader"
	"pvzService/pkg/responser"
)

type Params struct {
	fx.In

	Logger  *slog.Logger
	Usecase pvz.Usecase
}

type Handler struct {
	logger  *slog.Logger
	usecase pvz.Usecase
}

func NewHandler(params Params) *Handler {
	return &Handler{
		logger:  params.Logger,
		usecase: params.Usecase,
	}
}

func (h *Handler) CreatePvz(w http.ResponseWriter, r *http.Request) {
	pvzData := &models.PVZ{}
	if err := reader.ReadRequestData(r, pvzData); err != nil {
		h.logger.Error("create pvz request err: " + err.Error())
		responser.SendErr(w, http.StatusBadRequest, "incorrect data")
		return
	}

	createdPvz, err := h.usecase.CreatePvz(r.Context(), pvzData)
	if err != nil {
		// TODO: switch errors
		responser.SendErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	responser.SendOk(w, http.StatusCreated, createdPvz)
}

func (h *Handler) GetPvzList(w http.ResponseWriter, r *http.Request)         {}
func (h *Handler) CloseLastReception(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) DeleteLastProduct(w http.ResponseWriter, r *http.Request)  {}

func (h *Handler) CreateReception(w http.ResponseWriter, r *http.Request) {
	receptionData := &models.ReceptionRequest{}
	if err := reader.ReadRequestData(r, receptionData); err != nil {
		h.logger.Error("create reception request err: " + err.Error())
		responser.SendErr(w, http.StatusBadRequest, "incorrect data")
		return
	}
	createdReception, err := h.usecase.CreateReception(r.Context(), receptionData)
	if err != nil {
		// TODO: switch error
		responser.SendErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	responser.SendOk(w, http.StatusCreated, createdReception)
}

func (h *Handler) AddProduct(w http.ResponseWriter, r *http.Request) {}
