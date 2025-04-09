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
	"pvzService/internal/services/auth"
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

func (repo *Repository) GetUserByEmail(ctx context.Context, email string) (*models.Users, error) {
	query, args, err := repo.builder.
		Select("id", "email", "role", "password").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		repo.log.Error("build query error: " + err.Error())
		return nil, err
	}

	user := &models.Users{}
	err = repo.pool.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.Password,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			repo.log.Debug("get user by email: user not found")
			return nil, auth.ErrUserNotFound
		}
		repo.log.Error("failed to get user: " + err.Error())
		return nil, err
	}

	return user, nil

}
func (repo *Repository) CreateUser(ctx context.Context, user *models.Users) (*models.Users, error) {
	query, args, err := repo.builder.
		Insert("users").
		Columns("id", "email", "role", "password").
		Values(user.ID, user.Email, user.Role, user.Password).
		Suffix("RETURNING id, email, role, password").
		ToSql()
	if err != nil {
		repo.log.Error("build query error: " + err.Error())
		return nil, err
	}

	createdUser := &models.Users{}
	if err = repo.pool.QueryRow(ctx, query, args...).Scan(
		&createdUser.ID,
		&createdUser.Email,
		&createdUser.Role,
		&createdUser.Password,
	); err != nil {
		pgErr := &pgconn.PgError{}
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			repo.log.Warn("user already exists")
			return nil, auth.ErrAlreadyExists
		}
		repo.log.Error("failed to create user" + err.Error())
		return nil, err
	}
	return createdUser, nil
}
