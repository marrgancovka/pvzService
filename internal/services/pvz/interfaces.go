package pvz

import (
	"context"
	"github.com/google/uuid"
	"pvzService/internal/models"
)

type Usecase interface {
	CreatePvz(ctx context.Context, pvzData *models.PVZ) (*models.PVZ, error)
	CreateReception(ctx context.Context, receptionData *models.ReceptionRequest) (*models.Reception, error)
	CloseLastReceptions(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error)
	AddProduct(ctx context.Context, product *models.ProductRequest) (*models.Product, error)
	DeleteLastProduct(ctx context.Context, pvzId uuid.UUID) error
}

type Repository interface {
	CreatePvz(ctx context.Context, pvzData *models.PVZ) (*models.PVZ, error)
	CreateReception(ctx context.Context, receptionData *models.Reception) (*models.Reception, error)
	CloseLastReceptions(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error)
	AddProduct(ctx context.Context, product *models.Product, pvzID uuid.UUID) (*models.Product, error)
	DeleteLastProduct(ctx context.Context, pvzId uuid.UUID) error
}
