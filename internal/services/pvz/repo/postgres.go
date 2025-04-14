package repo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marrgancovka/pvzService/internal/models"
	"github.com/marrgancovka/pvzService/internal/services/pvz"
	"go.uber.org/fx"
	"log/slog"
	"time"
)

const (
	PgErrCodeAlreadyExists            = "23505"
	PgErrViolatesForeignKeyConstraint = "23503"
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

func (repo *Repository) CreatePvz(ctx context.Context, pvzData *models.Pvz) (*models.Pvz, error) {
	const op = "pvz.Repository.CreatePvz"
	logger := repo.log.With("op", op)

	query, args, err := repo.builder.
		Insert("pvz").
		Columns("id", "registration_date", "city").
		Values(pvzData.ID, pvzData.RegistrationDate, pvzData.City).
		Suffix("RETURNING id, registration_date, city").
		ToSql()
	if err != nil {
		logger.Error("build query error: " + err.Error())
		return nil, err
	}

	createdPvz := &models.Pvz{}
	if err = repo.pool.QueryRow(ctx, query, args...).Scan(
		&createdPvz.ID,
		&createdPvz.RegistrationDate,
		&createdPvz.City,
	); err != nil {
		pgErr := &pgconn.PgError{}
		if errors.As(err, &pgErr) && pgErr.Code == PgErrCodeAlreadyExists {
			logger.Warn("pvz with this id already exists")
			return nil, pvz.ErrAlreadyExists
		}
		logger.Error("failed to create pvz" + err.Error())
		return nil, err
	}

	return createdPvz, nil
}

func (repo *Repository) CreateReception(ctx context.Context, receptionData *models.Reception) (*models.Reception, error) {
	const op = "pvz.Repository.CreateReception"
	logger := repo.log.With("op", op)

	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		logger.Error("failed to begin transaction: " + err.Error())
		return nil, err
	}
	defer tx.Rollback(ctx)

	_, err = repo.getLastInProgressReceptionID(ctx, tx, receptionData.PvzID)
	if err == nil {
		logger.Error("failed to get last in-progress reception id: " + err.Error())
		return nil, pvz.ErrNoClosedReception
	}

	query, args, err := repo.builder.
		Insert("receptions").
		Columns("id", "date_time", "pvz_id", "status").
		Values(receptionData.ID, receptionData.DateTime, receptionData.PvzID, receptionData.Status).
		Suffix("RETURNING id, date_time, pvz_id, status").
		ToSql()
	if err != nil {
		logger.Error("build query error: " + err.Error())
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
		if errors.As(err, &pgErr) && pgErr.Code == PgErrViolatesForeignKeyConstraint {
			logger.Error("pvz with this id not exists")
			return nil, pvz.ErrPvzNotExists
		}
		logger.Error("failed to create reception: " + err.Error())
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		logger.Error("failed to commit transaction: " + err.Error())
		return nil, err
	}

	return createdReception, nil
}

func (repo *Repository) CloseLastReceptions(ctx context.Context, pvzId uuid.UUID) (*models.Reception, error) {
	const op = "pvz.Repository.CloseLastReceptions"
	logger := repo.log.With("op", op)

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
		logger.Error("build query error: " + err.Error())
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
			logger.Error("in this pvz open reception not exists")
			return nil, pvz.ErrNoOpenReception
		}
		logger.Error("failed to update status: " + err.Error())
		return nil, err
	}

	return closedReception, nil
}

func (repo *Repository) AddProduct(ctx context.Context, product *models.Product, pvzID uuid.UUID) (*models.Product, error) {
	const op = "pvz.Repository.AddProduct"
	logger := repo.log.With("op", op)

	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		logger.Error("failed to begin transaction: " + err.Error())
		return nil, err
	}
	defer tx.Rollback(ctx)

	inProgressReceptionID, err := repo.getLastInProgressReceptionID(ctx, tx, pvzID)
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
		logger.Error("build query error: " + err.Error())
		return nil, err
	}

	if err = tx.QueryRow(ctx, query, args...).Scan(
		&product.ID,
		&product.DateTime,
		&product.Type,
		&product.ReceptionID,
	); err != nil {
		logger.Error("failed to add product: " + err.Error())
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		logger.Error("failed to commit transaction: " + err.Error())
		return nil, err
	}

	return product, nil
}

func (repo *Repository) DeleteLastProduct(ctx context.Context, pvzID uuid.UUID) error {
	const op = "pvz.Repository.DeleteLastProduct"
	logger := repo.log.With("op", op)

	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		logger.Error("failed to begin transaction: " + err.Error())
		return err
	}
	defer tx.Rollback(ctx)

	inProgressReceptionID, err := repo.getLastInProgressReceptionID(ctx, tx, pvzID)
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
		logger.Error("build query error: " + err.Error())
		return err
	}

	productId := uuid.Nil
	if err = tx.QueryRow(ctx, query, args...).Scan(&productId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("no products in this reception")
			return pvz.ErrNoProduct
		}
		logger.Error("failed to found product: " + err.Error())
		return err
	}

	query, args, err = repo.builder.
		Delete("products").
		Where(squirrel.Eq{"id": productId}).
		ToSql()
	if err != nil {
		logger.Error("build query error: " + err.Error())
		return err
	}

	if _, err = tx.Exec(ctx, query, args...); err != nil {
		logger.Error("failed to delete product: " + err.Error())
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		logger.Error("failed to commit transaction: " + err.Error())
		return err
	}

	return nil
}

func (repo *Repository) getLastInProgressReceptionID(ctx context.Context, tx pgx.Tx, pvzID uuid.UUID) (uuid.UUID, error) {
	const op = "pvz.Repository.getLastInProgressReceptionID"
	logger := repo.log.With("op", op)

	query, args, err := repo.builder.
		Select("id").
		From("receptions").
		Where(squirrel.And{
			squirrel.Eq{"pvz_id": pvzID},
			squirrel.Eq{"status": models.StatusInProgress},
		}).
		ToSql()
	if err != nil {
		logger.Error("failed to build query: " + err.Error())
		return uuid.Nil, err
	}

	var receptionID uuid.UUID
	if err = tx.QueryRow(ctx, query, args...).Scan(&receptionID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn("no open reception found", "pvz_id", pvzID)
			return uuid.Nil, pvz.ErrNoOpenReception
		}
		logger.Error("failed to get open reception: " + err.Error())
		return uuid.Nil, err
	}

	return receptionID, nil
}

func (repo *Repository) GetPvz(ctx context.Context, startDate, endDate time.Time, limit, page uint64) ([]*models.PvzWithReceptions, error) {
	offset := (page - 1) * limit
	query := `SELECT
    p.id,
    p.registration_date,
    p.city,
    COALESCE(
            (
                SELECT json_agg(
                               json_build_object(
                                       'reception', json_build_object(
                                       'id', r.id,
                                       'dateTime', r.date_time,
                                       'pvzId', r.pvz_id,
                                       'status', r.status
                                                    ),
                                       'products', COALESCE((
                                           SELECT json_agg(
                                                          json_build_object(
                                                                  'id', pr.id,
                                                                  'dateTime', pr.date_time,
                                                                  'type', pr.type,
                                                                  'receptionId', pr.reception_id
                                                          )
                                                  )
                                           FROM products pr
                                           WHERE pr.reception_id = r.id
                                       ), '[]'::json)
                               )
                       )
                FROM receptions r
                WHERE r.pvz_id = p.id
                  AND r.date_time BETWEEN $1 AND $2
            ),
            '[]'::json
    ) AS receptions_json
FROM pvz p
ORDER BY registration_date DESC
LIMIT $3 OFFSET $4;`

	rows, err := repo.pool.Query(ctx, query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.PvzWithReceptions
	for rows.Next() {
		pvzWithReceptions := &models.PvzWithReceptions{}
		var receptionsJSON []byte

		err = rows.Scan(
			&pvzWithReceptions.ID,
			&pvzWithReceptions.RegistrationDate,
			&pvzWithReceptions.City,
			&receptionsJSON,
		)
		if err != nil {
			return nil, err
		}

		var receptionWithProducts []*models.ReceptionWithProducts
		if err = json.Unmarshal(receptionsJSON, &receptionWithProducts); err != nil {
			return nil, err
		}

		pvzWithReceptions.Receptions = receptionWithProducts

		result = append(result, pvzWithReceptions)
	}
	return result, nil
}

func (repo *Repository) GetPvzList(ctx context.Context) ([]*models.Pvz, error) {
	const op = "pvz.Repository.GetPvzList"
	logger := repo.log.With("op", op)

	query, _, err := repo.builder.
		Select("id", "registration_date", "city").
		From("pvz").
		OrderBy("registration_date").
		ToSql()
	if err != nil {
		logger.Error("failed to build query: " + err.Error())
		return nil, err
	}

	rows, err := repo.pool.Query(ctx, query)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		logger.Error("failed to execute query: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var results []*models.Pvz
	for rows.Next() {
		row := &models.Pvz{}
		if err = rows.Scan(&row.ID, &row.RegistrationDate, &row.City); err != nil {
			logger.Error("failed to scan row: " + err.Error())
			return nil, err
		}
		results = append(results, row)
	}

	return results, nil
}
