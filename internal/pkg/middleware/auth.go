package middleware

import (
	"context"
	"errors"
	"github.com/marrgancovka/pvzService/internal/pkg/jwter"
	"github.com/marrgancovka/pvzService/internal/services/auth"
	"github.com/marrgancovka/pvzService/pkg/responser"
	"go.uber.org/fx"
	"net/http"
	"strings"
)

const RoleInContext string = "RoleInContext"

type AuthMiddlewareParams struct {
	fx.In

	JWTer auth.JWTer
}

type AuthMiddleware struct {
	jwt auth.JWTer
}

func NewAuthMiddleware(p AuthMiddlewareParams) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: p.JWTer,
	}
}

func (authMD *AuthMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			responser.SendErr(w, http.StatusForbidden, "not authorized")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			responser.SendErr(w, http.StatusForbidden, "Authorization header format must be Bearer {token}")
			return
		}

		token := parts[1]

		claims, err := authMD.jwt.ValidateJWT(token)
		if err != nil && !errors.Is(err, jwter.ErrNoID) {
			responser.SendErr(w, http.StatusForbidden, "Invalid token")
			return
		}

		if !claims.Role.IsValid() {
			responser.SendErr(w, http.StatusForbidden, "Invalid role")
			return
		}

		ctx := context.WithValue(r.Context(), RoleInContext, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
