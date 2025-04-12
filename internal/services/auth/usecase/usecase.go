package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/marrgancovka/pvzService/internal/models"
	"github.com/marrgancovka/pvzService/internal/pkg/jwter"
	"github.com/marrgancovka/pvzService/internal/services/auth"
	"github.com/marrgancovka/pvzService/pkg/hasher"
	"go.uber.org/fx"
	"log/slog"
)

type Params struct {
	fx.In

	Logger *slog.Logger
	Repo   auth.Repository
	JWTer  auth.JWTer
}

type Usecase struct {
	log  *slog.Logger
	repo auth.Repository
	jwt  auth.JWTer
}

func NewUsecase(p Params) *Usecase {
	return &Usecase{
		log:  p.Logger,
		repo: p.Repo,
		jwt:  p.JWTer,
	}
}

func (uc *Usecase) DummyLogin(ctx context.Context, role *models.DummyLogin) (string, error) {
	if !role.Role.IsValid() {
		uc.log.Error("invalid role: " + string(role.Role))
		return "", auth.ErrIncorrectRole
	}

	tokenPayload := &models.TokenPayload{Role: role.Role}

	token, err := uc.jwt.GenerateJWT(tokenPayload)
	if err != nil && !errors.Is(err, jwter.ErrNoID) {
		return "", err
	}
	return token.Token, nil
}

func (uc *Usecase) Login(ctx context.Context, userData *models.Users) (string, error) {
	user, err := uc.repo.GetUserByEmail(ctx, userData.Email)
	if err != nil {
		return "", err
	}
	if !hasher.CompareStringHash(userData.Password, user.Password) {
		uc.log.Error("passwords don't match")
		return "", auth.ErrIncorrectData
	}

	tokenPayload := &models.TokenPayload{
		ID:   user.ID,
		Role: user.Role,
	}
	token, err := uc.jwt.GenerateJWT(tokenPayload)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

func (uc *Usecase) Register(ctx context.Context, userData *models.Users) (string, error) {
	if !userData.Role.IsValid() {
		uc.log.Error("invalid role: " + string(userData.Role))
		return "", auth.ErrIncorrectRole
	}
	userData.ID = uuid.New()
	userData.Password = hasher.GenerateHashString(userData.Password)
	newUser, err := uc.repo.CreateUser(ctx, userData)
	if err != nil {
		return "", err
	}
	tokenPayload := &models.TokenPayload{
		ID:   newUser.ID,
		Role: newUser.Role,
	}
	token, err := uc.jwt.GenerateJWT(tokenPayload)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}
