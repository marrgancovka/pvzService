package pvz

import (
	"context"
	"pvzService/internal/models"
)

type Usecase interface {
	CreatePvz(ctx context.Context, pvzData *models.PVZ) (*models.PVZ, error)
	CreateReception(ctx context.Context, receptionData *models.ReceptionRequest) (*models.Reception, error)
}

type Repository interface {
	CreatePvz(ctx context.Context, pvzData *models.PVZ) (*models.PVZ, error)
	CreateReception(ctx context.Context, receptionData *models.Reception) (*models.Reception, error)
}
