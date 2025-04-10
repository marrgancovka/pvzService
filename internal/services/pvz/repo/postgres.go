package repo

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
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
		return nil, pvz.ErrOpenReception
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
