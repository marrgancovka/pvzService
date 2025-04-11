package repo

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"log/slog"
	"pvzService/internal/models"
	"pvzService/internal/services/pvz"
)

type Params struct {
	fx.In

	Pool    *pgxpool.Pool
	Logger  *slog.Logger
	Builder squirrel.StatementBuilderType
}

type Repository struct {
	pool    *pgxpool.Pool
	log     *slog.Logger
	builder squirrel.StatementBuilderType
}

func NewRepository(params Params) *Repository {
	return &Repository{
		pool:    params.Pool,
		log:     params.Logger,
		builder: params.Builder,
	}
}

func (repo *Repository) CreatePvz(ctx context.Context, pvzData *models.PVZ) (*models.PVZ, error) {
	query, args, err := repo.builder.
		Insert("pvz").
		Columns("id", "registration_date", "city").
		Values(pvzData.ID, pvzData.RegistrationDate, pvzData.City).
		Suffix("RETURNING id, registration_date, city").
		ToSql()
	if err != nil {
		repo.log.Error("build query error: " + err.Error())
		return nil, err
	}

	createdPvz := &models.PVZ{}
	if err = repo.pool.QueryRow(ctx, query, args...).Scan(
		&createdPvz.ID,
		&createdPvz.RegistrationDate,
		&createdPvz.City,
	); err != nil {
		pgErr := &pgconn.PgError{}
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			repo.log.Warn("pvz with this id already exists")
			return nil, pvz.ErrAlreadyExists
		}
		repo.log.Error("failed to create pvz" + err.Error())
		return nil, err
	}
	return createdPvz, nil
}

// TODO: исправить?
func (repo *Repository) CreateReception(ctx context.Context, receptionData *models.Reception) (*models.Reception, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		repo.log.Error("failed to begin transaction: " + err.Error())
		return nil, err
	}
	defer tx.Rollback(ctx)

	query, args, err := repo.builder.
		Select("1").
		From("receptions").
		Where(squirrel.And{squirrel.Eq{"pvz_id": receptionData.PvzID}, squirrel.Eq{"status": models.StatusInProgress}}).
		ToSql()
	if err != nil {
		repo.log.Error("build query error: " + err.Error())
		return nil, err
	}

	haveInProgress := 0
	err = tx.QueryRow(ctx, query, args...).Scan(&haveInProgress)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		repo.log.Error("failed to check if reception exists: " + err.Error())
		return nil, err
	}
	if haveInProgress == 1 {
		repo.log.Warn("in this pvz open reception already exists")
		return nil, pvz.ErrNoOpenReception
	}

	query, args, err = repo.builder.
		Insert("receptions").
		Columns("id", "date_time", "pvz_id", "status").
		Values(receptionData.ID, receptionData.DateTime, receptionData.PvzID, receptionData.Status).
		Suffix("RETURNING id, date_time, pvz_id, status").
		ToSql()
	if err != nil {
		repo.log.Error("build query error: " + err.Error())
		return nil, err
	}

	createdReception := &models.Reception{}
	if err = tx.QueryRow(ctx, query, args...).Scan(
		&createdReception.ID,
		&createdReception.DateTime,
		&createdReception.PvzID,
		&createdReception.Status,
	); err != nil {
		pgErr := &pgconn.PgError{}
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			repo.log.Warn("pvz with this id not exists")
			return nil, pvz.ErrPvzNotExists
		}
		repo.log.Error("failed to create reception: " + err.Error())
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		repo.log.Error("failed to commit transaction: " + err.Error())
		return nil, err
	}

	return createdReception, nil
}

func (repo *Repository) CloseLastReceptions(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error) {
	query, args, err := repo.builder.
		Update("receptions").
		Set("status", models.StatusClose).
		Where(squirrel.And{
			squirrel.Eq{"pvz_id": pvzId},
			squirrel.Eq{"status": models.StatusInProgress},
		}).
		Suffix("RETURNING id, date_time, pvz_id, status").
		ToSql()
	if err != nil {
		repo.log.Error("build query error: " + err.Error())
		return nil, err
	}

	closedReception := &models.Reception{}
	if err = repo.pool.QueryRow(ctx, query, args...).Scan(
		&closedReception.ID,
		&closedReception.DateTime,
		&closedReception.PvzID,
		&closedReception.Status,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			repo.log.Warn("in this pvz open reception not exists")
			return nil, pvz.ErrNoClosedReception
		}
		repo.log.Error("failed to update status: " + err.Error())
		return nil, err
	}

	return closedReception, nil
}

func (repo *Repository) AddProduct(ctx context.Context, product *models.Product, pvzID uuid.UUID) (*models.Product, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		repo.log.Error("failed to begin transaction: " + err.Error())
		return nil, err
	}
	defer tx.Rollback(ctx)

	inProgressReceptionID, err := repo.GetLastInProgressReceptionID(ctx, tx, pvzID)
	if err != nil {
		return nil, err
	}

	query, args, err := repo.builder.
		Insert("products").
		Columns("id", "date_time", "type", "reception_id").
		Values(product.ID, product.DateTime, product.Type, inProgressReceptionID).
		Suffix("RETURNING id, date_time, type, reception_id").
		ToSql()
	if err != nil {
		repo.log.Error("build query error: " + err.Error())
		return nil, err
	}

	if err = tx.QueryRow(ctx, query, args...).Scan(
		&product.ID,
		&product.DateTime,
		&product.Type,
		&product.ReceptionID,
	); err != nil {
		// TODO: errors?
		repo.log.Error("failed to add product: " + err.Error())
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		repo.log.Error("failed to commit transaction: " + err.Error())
		return nil, err
	}

	return product, nil
}

func (repo *Repository) DeleteLastProduct(ctx context.Context, pvzID uuid.UUID) error {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		repo.log.Error("failed to begin transaction: " + err.Error())
	}
	defer tx.Rollback(ctx)

	inProgressReceptionID, err := repo.GetLastInProgressReceptionID(ctx, tx, pvzID)
	if err != nil {
		return err
	}

	query, args, err := repo.builder.
		Select("id").
		From("products").
		Where(squirrel.Eq{"reception_id": inProgressReceptionID}).
		OrderBy("date_time DESC").
		Limit(1).
		ToSql()
	if err != nil {
		repo.log.Error("build query error: " + err.Error())
		return err
	}

	productId := uuid.Nil
	if err = tx.QueryRow(ctx, query, args...).Scan(&productId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			repo.log.Warn("нет продуктов в этой приемке")
			return pvz.ErrNoProduct
		}
		repo.log.Error("failed to found product: " + err.Error())
		return err
	}

	query, args, err = repo.builder.
		Delete("products").
		Where(squirrel.Eq{"id": productId}).
		ToSql()
	if err != nil {
		repo.log.Error("build query error: " + err.Error())
		return err
	}

	if _, err = tx.Exec(ctx, query, args...); err != nil {
		repo.log.Error("failed to delete product: " + err.Error())
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		repo.log.Error("failed to commit transaction: " + err.Error())
		return err
	}
	return nil
}

func (repo *Repository) GetLastInProgressReceptionID(ctx context.Context, tx pgx.Tx, pvzID uuid.UUID) (uuid.UUID, error) {
	query, args, err := repo.builder.
		Select("id").
		From("receptions").
		Where(squirrel.And{
			squirrel.Eq{"pvz_id": pvzID},
			squirrel.Eq{"status": models.StatusInProgress},
		}).
		ToSql()
	if err != nil {
		repo.log.Error("failed to build query", "error", err)
		return uuid.Nil, err
	}

	var receptionID uuid.UUID
	if err = tx.QueryRow(ctx, query, args...).Scan(&receptionID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			repo.log.Warn("no open reception found", "pvz_id", pvzID)
			return uuid.Nil, pvz.ErrNoOpenReception
		}
		repo.log.Error("failed to get open reception", "error", err)
		return uuid.Nil, err
	}

	return receptionID, nil
}
