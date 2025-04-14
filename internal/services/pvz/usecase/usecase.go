package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/marrgancovka/pvzService/internal/models"
	"github.com/marrgancovka/pvzService/internal/services/pvz"
	"go.uber.org/fx"
	"log/slog"
	"time"
)

type Params struct {
	fx.In

	Logger *slog.Logger
	Repo   pvz.Repository
}

type Usecase struct {
	log  *slog.Logger
	repo pvz.Repository
}

func NewUsecase(p Params) *Usecase {
	return &Usecase{
		log:  p.Logger,
		repo: p.Repo,
	}
}

func (uc *Usecase) CreatePvz(ctx context.Context, pvzData *models.Pvz) (*models.Pvz, error) {
	const op = "pvz.Usecase.CreatePvz"
	logger := uc.log.With("op", op)

	if !pvzData.City.IsValid() {
		logger.Error("incorrect city: " + string(pvzData.City))
		return nil, pvz.ErrInaccessibleCity
	}

	if pvzData.ID == uuid.Nil {
		logger.Warn("id in pvz data is nil, uuid generated")
		pvzData.ID = uuid.New()
	}

	nilTime := time.Time{}
	if pvzData.RegistrationDate == nilTime {
		logger.Warn("registration date in pvz data is nil, set now")
		pvzData.RegistrationDate = time.Now()
	}

	createdPvz, err := uc.repo.CreatePvz(ctx, pvzData)
	if err != nil {
		return nil, err
	}
	return createdPvz, nil
}

func (uc *Usecase) CreateReception(ctx context.Context, receptionData *models.ReceptionRequest) (*models.Reception, error) {
	const op = "pvz.Usecase.CreateReception"

	reception := &models.Reception{
		ID:       uuid.New(),
		DateTime: time.Now(),
		PvzID:    receptionData.PvzID,
		Status:   models.StatusInProgress,
	}

	createdReception, err := uc.repo.CreateReception(ctx, reception)
	if err != nil {
		return nil, err
	}

	return createdReception, nil
}

func (uc *Usecase) CloseLastReceptions(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error) {
	const op = "pvz.Usecase.CloseLastReceptions"

	closedReception, err := uc.repo.CloseLastReceptions(ctx, pvzId)
	if err != nil {
		return nil, err
	}
	return closedReception, nil
}

func (uc *Usecase) AddProduct(ctx context.Context, product *models.ProductRequest) (*models.Product, error) {
	const op = "pvz.Usecase.AddProduct"
	logger := uc.log.With("op", op)

	if !product.Type.IsValid() {
		logger.Error("incorrect type for product: " + string(product.Type))
		return nil, pvz.ErrIncorrectProductType
	}

	productData := &models.Product{
		ID:          uuid.New(),
		DateTime:    time.Now(),
		Type:        product.Type,
		ReceptionID: product.PvzID,
	}

	addedProduct, err := uc.repo.AddProduct(ctx, productData, product.PvzID)
	if err != nil {
		return nil, err
	}

	return addedProduct, nil
}

func (uc *Usecase) DeleteLastProduct(ctx context.Context, pvzId uuid.UUID) error {
	const op = "pvz.Usecase.DeleteLastProduct"

	err := uc.repo.DeleteLastProduct(ctx, pvzId)
	if err != nil {
		return err
	}
	return nil
}

func (uc *Usecase) GetPvz(ctx context.Context, startDate, endDate time.Time, limit, page uint64) (*models.Pvz, error) {
	const op = "pvz.Usecase.GetPvz"

	return uc.repo.GetPvz(ctx, startDate, endDate, limit, page)
}

func (uc *Usecase) GetPvzList(ctx context.Context) ([]*models.Pvz, error) {
	const op = "pvz.Usecase.GetPvzList"

	return uc.repo.GetPvzList(ctx)
}
