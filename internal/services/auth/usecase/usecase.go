package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"log/slog"
	"pvzService/internal/models"
	"pvzService/internal/pkg/jwter"
	"pvzService/internal/services/auth"
	"pvzService/pkg/hasher"
)

type Params struct {
	fx.In

	Logger *slog.Logger
	Repo   auth.Repository
	JWTer  *jwter.JWTer
}

type Usecase struct {
	log   *slog.Logger
	repo  auth.Repository
	JWTer *jwter.JWTer
}

func NewUsecase(p Params) *Usecase {
	return &Usecase{
		log:   p.Logger,
		repo:  p.Repo,
		JWTer: p.JWTer,
	}
}

func (uc *Usecase) DummyLogin(ctx context.Context, role *models.DummyLogin) (string, error) {
	tokenPayload := &models.TokenPayload{Role: role.Role}

	token, err := uc.JWTer.GenerateJWT(tokenPayload)
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
		uc.log.Error("invalid password")
		return "", auth.ErrIncorrectData
	}

	tokenPayload := &models.TokenPayload{
		ID:   user.ID,
		Role: user.Role,
	}
	token, err := uc.JWTer.GenerateJWT(tokenPayload)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}
func (uc *Usecase) Register(ctx context.Context, userData *models.Users) (string, error) {
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
	token, err := uc.JWTer.GenerateJWT(tokenPayload)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}
