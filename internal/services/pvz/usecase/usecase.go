package usecase

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"log/slog"
	"pvzService/internal/models"
	"pvzService/internal/pkg/jwter"
	"pvzService/internal/services/pvz"
	"time"
)

type Params struct {
	fx.In

	Logger *slog.Logger
	Repo   pvz.Repository
	JWTer  *jwter.JWTer
}

type Usecase struct {
	log   *slog.Logger
	repo  pvz.Repository
	JWTer *jwter.JWTer
}

func NewUsecase(p Params) *Usecase {
	return &Usecase{
		log:   p.Logger,
		repo:  p.Repo,
		JWTer: p.JWTer,
	}
}

func (uc *Usecase) CreatePvz(ctx context.Context, pvzData *models.Pvz) (*models.Pvz, error) {
	if !pvzData.City.IsValid() {
		uc.log.Error("incorrect city: " + string(pvzData.City))
		return nil, pvz.ErrInaccessibleCity
	}
	if pvzData.ID == uuid.Nil {
		uc.log.Warn("create pvz: id is nil")
		pvzData.ID = uuid.New()
	}
	nilTime := time.Time{}
	if pvzData.RegistrationDate == nilTime {
		uc.log.Warn("create pvz: registration date is nil")
		pvzData.RegistrationDate = time.Now()
	}

	createdPvz, err := uc.repo.CreatePvz(ctx, pvzData)
	if err != nil {
		return nil, err
	}
	return createdPvz, nil
}

func (uc *Usecase) CreateReception(ctx context.Context, receptionData *models.ReceptionRequest) (*models.Reception, error) {
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
	closedReception, err := uc.repo.CloseLastReceptions(ctx, pvzId)
	if err != nil {
		return nil, err
	}
	return closedReception, nil
}

func (uc *Usecase) AddProduct(ctx context.Context, product *models.ProductRequest) (*models.Product, error) {
	if !product.Type.IsValid() {
		uc.log.Error("incorrect type for product: " + string(product.Type))
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
	err := uc.repo.DeleteLastProduct(ctx, pvzId)
	if err != nil {
		return err
	}
	return nil
}

func (uc *Usecase) GetPvz(ctx context.Context, startDate, endDate time.Time, limit, page uint64) (*models.Pvz, error) {
	return uc.repo.GetPvz(ctx, startDate, endDate, limit, page)
}
