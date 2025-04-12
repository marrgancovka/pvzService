package auth

import (
	"context"
	"github.com/marrgancovka/pvzService/internal/models"
)

type Usecase interface {
	DummyLogin(ctx context.Context, role *models.DummyLogin) (string, error)
	Login(ctx context.Context, userData *models.Users) (string, error)
	Register(ctx context.Context, userData *models.Users) (string, error)
}

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.Users, error)
	CreateUser(ctx context.Context, user *models.Users) (*models.Users, error)
}
