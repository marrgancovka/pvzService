package jwter

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/marrgancovka/pvzService/internal/models"
	"go.uber.org/fx"
	"log/slog"
	"time"
)

type Params struct {
	fx.In

	Config Config
	Logger *slog.Logger
}

type JWTer struct {
	cfg Config
	log *slog.Logger
}

func New(p Params) *JWTer {
	return &JWTer{
		cfg: p.Config,
		log: p.Logger,
	}
}

func (jwter *JWTer) GenerateJWT(payload *models.TokenPayload) (*models.Token, error) {
	expTime := time.Now().Add(jwter.cfg.ExpirationTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub":  payload.ID,
		"role": payload.Role,
		"exp":  expTime.Unix(),
	})

	tokenStr, err := token.SignedString(jwter.cfg.KeyJWT)
	if err != nil {
		jwter.log.Error("JWT Error: " + err.Error())
		return nil, err
	}

	tokenResponse := &models.Token{
		Token: tokenStr,
	}

	if payload.ID == uuid.Nil {
		jwter.log.Warn("no id in payload")
		err = ErrNoID
	}

	return tokenResponse, err
}

func (jwter *JWTer) ValidateJWT(tokenString string) (*models.TokenPayload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			jwter.log.Error("validate jwt: Unexpected signing method")
			return nil, ErrUnexpectedSigningMethod
		}

		return jwter.cfg.KeyJWT, nil
	})
	if err != nil {
		jwter.log.Error("parsing token: " + err.Error())
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		jwter.log.Error("jwt validate: claims error")
		return nil, ErrInvalidToken
	}

	id, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		jwter.log.Error("jwt validate: invalid id")
		return nil, ErrInvalidTokenClaims
	}

	role, ok := claims["role"].(string)
	if !ok {
		jwter.log.Error("jwt validate: invalid role: " + claims["role"].(string))
		return nil, ErrInvalidTokenClaims
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		jwter.log.Error("jwt validate: invalid exp")
		return nil, ErrInvalidTokenClaims
	}
	expTime := time.Unix(int64(exp), 0)

	if expTime.Before(time.Now()) {
		jwter.log.Error("validate jwt: token expired")
		return nil, ErrTokenExpired
	}

	payload := &models.TokenPayload{
		ID:   id,
		Role: models.Role(role),
		Exp:  expTime,
	}

	if id == uuid.Nil {
		jwter.log.Warn("validate jwt: no id")
		err = ErrNoID
	}

	return payload, err
}
