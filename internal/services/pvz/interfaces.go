package pvz

import (
	"context"
	"github.com/google/uuid"
	"github.com/marrgancovka/pvzService/internal/models"
	"time"
)

type Usecase interface {
	CreatePvz(ctx context.Context, pvzData *models.Pvz) (*models.Pvz, error)
	CreateReception(ctx context.Context, receptionData *models.ReceptionRequest) (*models.Reception, error)
	CloseLastReceptions(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error)
	AddProduct(ctx context.Context, product *models.ProductRequest) (*models.Product, error)
	DeleteLastProduct(ctx context.Context, pvzId uuid.UUID) error
	GetPvz(ctx context.Context, startDate, endDate time.Time, limit, page uint64) (*models.Pvz, error)
	GetPvzList(ctx context.Context) ([]*models.Pvz, error)
}

type Repository interface {
	CreatePvz(ctx context.Context, pvzData *models.Pvz) (*models.Pvz, error)
	CreateReception(ctx context.Context, receptionData *models.Reception) (*models.Reception, error)
	CloseLastReceptions(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error)
	AddProduct(ctx context.Context, product *models.Product, pvzID uuid.UUID) (*models.Product, error)
	DeleteLastProduct(ctx context.Context, pvzId uuid.UUID) error
	GetPvz(ctx context.Context, startDate, endDate time.Time, limit, page uint64) (*models.Pvz, error)
	GetPvzList(ctx context.Context) ([]*models.Pvz, error)
}
