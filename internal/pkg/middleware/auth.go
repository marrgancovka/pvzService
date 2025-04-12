package middleware

import (
	"context"
	"errors"
	"go.uber.org/fx"
	"net/http"
	"pvzService/internal/pkg/jwter"
	"pvzService/pkg/responser"
	"strings"
)

const RoleInContext string = "RoleInContext"

type AuthMiddlewareParams struct {
	fx.In

	JWTer *jwter.JWTer
}

type AuthMiddleware struct {
	jwt *jwter.JWTer
}

func NewAuthMiddleware(jwter *jwter.JWTer) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: jwter,
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
