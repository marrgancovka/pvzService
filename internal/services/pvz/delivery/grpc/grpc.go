package grpc

import (
	"context"
	"github.com/marrgancovka/pvzService/internal/models"
	"github.com/marrgancovka/pvzService/internal/services/pvz"
	"github.com/marrgancovka/pvzService/internal/services/pvz/delivery/grpc/gen"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
)

type Params struct {
	fx.In

	Logger  *slog.Logger
	Usecase pvz.Usecase
}

type Handler struct {
	logger  *slog.Logger
	usecase pvz.Usecase

	gen.PVZServiceServer
}

func NewHandler(params Params) *Handler {
	return &Handler{
		logger:  params.Logger,
		usecase: params.Usecase,
	}
}

func (h *Handler) GetPVZList(ctx context.Context, _ *gen.GetPVZListRequest) (*gen.GetPVZListResponse, error) {
	const op = "grpc.pvz.Handler.GetPVZList"
	h.logger = h.logger.With("op", op)

	results, err := h.usecase.GetPvzList(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	converted := make([]*gen.PVZ, len(results))
	for i := range results {
		converted[i] = convert(results[i])
	}

	h.logger.Info("success get pvz list", "result", converted)
	return &gen.GetPVZListResponse{Pvzs: converted}, nil
}

func convert(pvz *models.Pvz) *gen.PVZ {
	return &gen.PVZ{
		Id:               pvz.ID.String(),
		RegistrationDate: timestamppb.New(pvz.RegistrationDate),
		City:             string(pvz.City),
	}
}
