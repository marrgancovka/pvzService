package http

import (
	"fmt"
	"go.uber.org/fx"
	"log/slog"
	"net/http"
	"pvzService/internal/models"
	"pvzService/internal/pkg/middleware"
	"pvzService/internal/services/pvz"
	"pvzService/pkg/reader"
	"pvzService/pkg/responser"
	"strconv"
	"time"
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
	if r.Context().Value(middleware.RoleInContext) != models.RoleModerator {
		h.logger.Error("create pvz: only for moderator")
		responser.SendErr(w, http.StatusForbidden, "нет доступа")
		return
	}

	pvzData := &models.Pvz{}
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

func (h *Handler) GetPvzList(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(middleware.RoleInContext) != models.RoleModerator && r.Context().Value(middleware.RoleInContext) != models.RoleEmployee {
		h.logger.Error("get list pvz: only for moderator or employee")
		responser.SendErr(w, http.StatusForbidden, "нет доступа")
		return
	}

	var err error
	queryParams := r.URL.Query()

	startDateStr := queryParams.Get("startDate")
	startDate := time.Now().AddDate(-1, 0, 0)
	if startDateStr != "" {
		startDate, err = parseDate(startDateStr)
		if err != nil {
			h.logger.Error("get pvz list request err: " + err.Error())
			responser.SendErr(w, http.StatusBadRequest, "incorrect data")
			return
		}
	}

	endDateStr := queryParams.Get("endDate")
	endDate := time.Now().AddDate(-1, 0, 0)
	if endDateStr != "" {
		endDate, err = parseDate(endDateStr)
		if err != nil {
			h.logger.Error("get pvz list request err: " + err.Error())
			responser.SendErr(w, http.StatusBadRequest, "incorrect data")
			return
		}
	}

	limitStr := queryParams.Get("limit")
	limit := uint64(10)
	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil {
			h.logger.Error("get pvz list request err: " + err.Error())
			responser.SendErr(w, http.StatusBadRequest, "incorrect data")
			return
		}
		limit = uint64(limitInt)
	}

	pageStr := queryParams.Get("page")
	page := uint64(1)
	if pageStr != "" {
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil {
			h.logger.Error("get pvz list request err: " + err.Error())
			responser.SendErr(w, http.StatusBadRequest, "incorrect data")
			return
		}
		page = uint64(pageInt)
	}

	pvzList, err := h.usecase.GetPvz(r.Context(), startDate, endDate, limit, page)
	if err != nil {
		// TODO: switch errors
		responser.SendErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	responser.SendOk(w, http.StatusOK, pvzList)
}

func (h *Handler) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(middleware.RoleInContext) != models.RoleEmployee {
		h.logger.Error("create pvz: only for employee")
		responser.SendErr(w, http.StatusForbidden, "нет доступа")
		return
	}

	pvzId, err := reader.ReadVarsUUID(r, "pvzId")
	if err != nil {
		h.logger.Error("get pvz id err: " + err.Error())
		responser.SendErr(w, http.StatusBadRequest, "incorrect data")
		return
	}

	closedReception, err := h.usecase.CloseLastReceptions(r.Context(), pvzId)
	if err != nil {
		// TODO: switch errors
		responser.SendErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	responser.SendOk(w, http.StatusOK, closedReception)
}

func (h *Handler) DeleteLastProduct(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(middleware.RoleInContext) != models.RoleEmployee {
		h.logger.Error("create pvz: only for employee")
		responser.SendErr(w, http.StatusForbidden, "нет доступа")
		return
	}

	pvzId, err := reader.ReadVarsUUID(r, "pvzId")
	if err != nil {
		h.logger.Error("get pvz id err: " + err.Error())
		responser.SendErr(w, http.StatusBadRequest, "incorrect id")
		return
	}

	err = h.usecase.DeleteLastProduct(r.Context(), pvzId)
	if err != nil {
		// TODO: switch errors
		responser.SendErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	responser.SendOk(w, http.StatusOK, "Товар удален")
}

func (h *Handler) CreateReception(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(middleware.RoleInContext) != models.RoleEmployee {
		h.logger.Error("create pvz: only for employee")
		responser.SendErr(w, http.StatusForbidden, "нет доступа")
		return
	}

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

func (h *Handler) AddProduct(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(middleware.RoleInContext) != models.RoleEmployee {
		h.logger.Error("create pvz: only for employee")
		responser.SendErr(w, http.StatusForbidden, "нет доступа")
		return
	}

	productData := &models.ProductRequest{}
	if err := reader.ReadRequestData(r, productData); err != nil {
		h.logger.Error("create product request err: " + err.Error())
		responser.SendErr(w, http.StatusBadRequest, "incorrect data")
		return
	}

	addedProduct, err := h.usecase.AddProduct(r.Context(), productData)
	if err != nil {
		// TODO: switch errors
		responser.SendErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	responser.SendOk(w, http.StatusCreated, addedProduct)
}

func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unrecognized date format")
}
