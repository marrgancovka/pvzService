package http

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/marrgancovka/pvzService/internal/models"
	"github.com/marrgancovka/pvzService/internal/pkg/metrics"
	"github.com/marrgancovka/pvzService/internal/pkg/middleware"
	"github.com/marrgancovka/pvzService/internal/services/pvz"
	"github.com/marrgancovka/pvzService/internal/services/pvz/delivery/grpc/gen"
	"github.com/marrgancovka/pvzService/pkg/reader"
	"github.com/marrgancovka/pvzService/pkg/responser"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Params struct {
	fx.In

	Logger     *slog.Logger
	Usecase    pvz.Usecase
	GRPCClient *grpc.ClientConn
	Metrics    metrics.Metrics
}

type Handler struct {
	logger     *slog.Logger
	usecase    pvz.Usecase
	grpcClient gen.PVZServiceClient
	metrics    metrics.Metrics
}

func NewHandler(params Params) *Handler {
	return &Handler{
		grpcClient: gen.NewPVZServiceClient(params.GRPCClient),
		logger:     params.Logger,
		usecase:    params.Usecase,
		metrics:    params.Metrics,
	}
}

func (h *Handler) CreatePvz(w http.ResponseWriter, r *http.Request) {
	const op = "pvz.Handler.CreatePvz"
	logger := h.logger.With("op", op)

	if r.Context().Value(middleware.RoleInContext) != models.RoleModerator {
		logger.Error("only for moderator")
		responser.SendErr(w, http.StatusForbidden, pvz.ErrNoAccess.Error())
		return
	}

	pvzData := &models.Pvz{}
	if err := reader.ReadRequestData(r, pvzData); err != nil {
		logger.Error("error read request data: " + err.Error())
		responser.SendErr(w, http.StatusBadRequest, pvz.ErrBadRequest.Error())
		return
	}

	createdPvz, err := h.usecase.CreatePvz(r.Context(), pvzData)
	if err != nil {
		switch {
		case errors.Is(err, pvz.ErrAlreadyExists):
			responser.SendErr(w, http.StatusBadRequest, pvz.ErrBadRequest.Error())
			return
		case errors.Is(err, pvz.ErrInaccessibleCity):
			responser.SendErr(w, http.StatusBadRequest, pvz.ErrInaccessibleCity.Error())
			return
		default:
			responser.SendErr(w, http.StatusInternalServerError, "internal mainServer error")
			return
		}
	}

	logger.Info("success create pvz", "response", createdPvz)
	h.metrics.CreatedPvzTotal(string(createdPvz.City))
	responser.SendOk(w, http.StatusCreated, createdPvz)
}

func (h *Handler) GetPvzs(w http.ResponseWriter, r *http.Request) {
	const op = "pvz.Handler.GetPvzList"
	logger := h.logger.With("op", op)

	if r.Context().Value(middleware.RoleInContext) != models.RoleModerator && r.Context().Value(middleware.RoleInContext) != models.RoleEmployee {
		logger.Error("only for moderator or employee")
		responser.SendErr(w, http.StatusForbidden, pvz.ErrNoAccess.Error())
		return
	}

	var err error
	queryParams := r.URL.Query()

	startDateStr := queryParams.Get("startDate")
	startDate := time.Now().AddDate(-1, 0, 0)
	if startDateStr != "" {
		startDate, err = parseDate(startDateStr)
		if err != nil {
			logger.Error("error parsing date: " + err.Error())
			responser.SendErr(w, http.StatusBadRequest, pvz.ErrBadRequest.Error())
			return
		}
	}

	endDateStr := queryParams.Get("endDate")
	endDate := time.Now().AddDate(-1, 0, 0)
	if endDateStr != "" {
		endDate, err = parseDate(endDateStr)
		if err != nil {
			logger.Error("error parsing date: " + err.Error())
			responser.SendErr(w, http.StatusBadRequest, pvz.ErrBadRequest.Error())
			return
		}
	}

	limitStr := queryParams.Get("limit")
	limit := uint64(10)
	if limitStr != "" {
		limitInt, err := strconv.Atoi(limitStr)
		if err != nil {
			logger.Error("error conv limit: " + err.Error())
			responser.SendErr(w, http.StatusBadRequest, pvz.ErrBadRequest.Error())
			return
		}
		limit = uint64(limitInt)
	}

	pageStr := queryParams.Get("page")
	page := uint64(1)
	if pageStr != "" {
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil {
			logger.Error("error conv page: " + err.Error())
			responser.SendErr(w, http.StatusBadRequest, pvz.ErrBadRequest.Error())
			return
		}
		page = uint64(pageInt)
	}

	pvzList, err := h.usecase.GetPvz(r.Context(), startDate, endDate, limit, page)
	if err != nil {
		responser.SendErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("success get pvz", "response", pvzList)
	responser.SendOk(w, http.StatusOK, pvzList)
}

func (h *Handler) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	const op = "pvz.Handler.CloseLastReception"
	logger := h.logger.With("op", op)

	if r.Context().Value(middleware.RoleInContext) != models.RoleEmployee {
		logger.Error("only for employee")
		responser.SendErr(w, http.StatusForbidden, pvz.ErrNoAccess.Error())
		return
	}

	pvzId, err := reader.ReadVarsUUID(r, "pvzId")
	if err != nil {
		logger.Error("error read var uuid: " + err.Error())
		responser.SendErr(w, http.StatusBadRequest, pvz.ErrBadRequest.Error())
		return
	}
	if pvzId == uuid.Nil {
		logger.Error("pvz id is nil")
		responser.SendErr(w, http.StatusBadRequest, pvz.ErrBadRequest.Error())
		return
	}

	closedReception, err := h.usecase.CloseLastReceptions(r.Context(), pvzId)
	if err != nil {
		switch {
		case errors.Is(err, pvz.ErrNoOpenReception):
			responser.SendErr(w, http.StatusBadRequest, pvz.ErrNoOpenReception.Error())
			return
		default:
			responser.SendErr(w, http.StatusInternalServerError, "internal mainServer error")
			return
		}
	}

	logger.Info("success close last reception", "response", closedReception)
	responser.SendOk(w, http.StatusOK, closedReception)
}

func (h *Handler) DeleteLastProduct(w http.ResponseWriter, r *http.Request) {
	const op = "pvz.Handler.DeleteLastProduct"
	logger := h.logger.With("op", op)

	if r.Context().Value(middleware.RoleInContext) != models.RoleEmployee {
		logger.Error("only for employee")
		responser.SendErr(w, http.StatusForbidden, pvz.ErrNoAccess.Error())
		return
	}

	pvzId, err := reader.ReadVarsUUID(r, "pvzId")
	if err != nil {
		logger.Error("error read var uuid: " + err.Error())
		responser.SendErr(w, http.StatusBadRequest, pvz.ErrBadRequest.Error())
		return
	}

	err = h.usecase.DeleteLastProduct(r.Context(), pvzId)
	if err != nil {
		switch {
		case errors.Is(err, pvz.ErrNoOpenReception):
			responser.SendErr(w, http.StatusBadRequest, pvz.ErrNoOpenReception.Error())
			return
		case errors.Is(err, pvz.ErrNoProduct):
			responser.SendErr(w, http.StatusBadRequest, pvz.ErrNoProduct.Error())
			return
		default:
			responser.SendErr(w, http.StatusInternalServerError, "internal mainServer error")
			return
		}
	}

	logger.Info("success delete last product")
	responser.SendOk(w, http.StatusOK, "Товар удален")
}

func (h *Handler) CreateReception(w http.ResponseWriter, r *http.Request) {
	op := "pvz.Handler.CreateReception"
	logger := h.logger.With("op", op)

	if r.Context().Value(middleware.RoleInContext) != models.RoleEmployee {
		logger.Error("only for employee")
		responser.SendErr(w, http.StatusForbidden, pvz.ErrNoAccess.Error())
		return
	}

	receptionData := &models.ReceptionRequest{}
	if err := reader.ReadRequestData(r, receptionData); err != nil {
		logger.Error("error read request data: " + err.Error())
		responser.SendErr(w, http.StatusBadRequest, pvz.ErrBadRequest.Error())
		return
	}
	createdReception, err := h.usecase.CreateReception(r.Context(), receptionData)
	if err != nil {
		switch {
		case errors.Is(err, pvz.ErrNoOpenReception):
			responser.SendErr(w, http.StatusBadRequest, pvz.ErrNoOpenReception.Error())
			return
		default:
			responser.SendErr(w, http.StatusInternalServerError, "internal mainServer error")
			return
		}
	}

	logger.Info("success create reception", "response", createdReception)
	h.metrics.CreatedReceptionsTotal(fmt.Sprint(createdReception.PvzID))
	responser.SendOk(w, http.StatusCreated, createdReception)
}

func (h *Handler) AddProduct(w http.ResponseWriter, r *http.Request) {
	op := "pvz.Handler.AddProduct"
	logger := h.logger.With("op", op)

	if r.Context().Value(middleware.RoleInContext) != models.RoleEmployee {
		logger.Error("only for employee")
		responser.SendErr(w, http.StatusForbidden, pvz.ErrNoAccess.Error())
		return
	}

	productData := &models.ProductRequest{}
	if err := reader.ReadRequestData(r, productData); err != nil {
		logger.Error("error read request data: " + err.Error())
		responser.SendErr(w, http.StatusBadRequest, pvz.ErrBadRequest.Error())
		return
	}

	if productData.PvzID == uuid.Nil {
		logger.Error("pvz id is nil")
		responser.SendErr(w, http.StatusBadRequest, pvz.ErrBadRequest.Error())
		return
	}

	addedProduct, err := h.usecase.AddProduct(r.Context(), productData)
	if err != nil {
		switch {
		case errors.Is(err, pvz.ErrNoOpenReception):
			responser.SendErr(w, http.StatusBadRequest, pvz.ErrNoOpenReception.Error())
			return
		default:
			responser.SendErr(w, http.StatusInternalServerError, "internal mainServer error")
			return
		}
	}

	logger.Info("success add product", "response", addedProduct)
	h.metrics.AddedProductTotal(string(addedProduct.Type))
	responser.SendOk(w, http.StatusCreated, addedProduct)
}

func (h *Handler) GetPvzList(w http.ResponseWriter, r *http.Request) {
	const op = "pvz.Handler.GetPvzList"
	logger := h.logger.With("op", op)

	list, err := h.grpcClient.GetPVZList(r.Context(), &gen.GetPVZListRequest{})
	if err != nil {
		logger.Error("error get pvz list: " + err.Error())
		responser.SendErr(w, http.StatusInternalServerError, "internal mainServer error")
	}

	result := make([]*models.Pvz, len(list.Pvzs))
	for i := range list.Pvzs {
		result[i] = convert(list.Pvzs[i])
	}

	logger.Info("success get pvz list on grpc", "response", result)
	responser.SendOk(w, http.StatusOK, result)
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

func convert(pvz *gen.PVZ) *models.Pvz {
	pvzUUID, err := uuid.Parse(pvz.Id)
	if err != nil {
		return nil
	}

	return &models.Pvz{
		ID:               pvzUUID,
		RegistrationDate: pvz.RegistrationDate.AsTime(),
		City:             models.City(pvz.City),
	}
}
