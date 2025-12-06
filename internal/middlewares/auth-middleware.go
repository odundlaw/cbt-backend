// Package middlewares where all middlware files are stored for proper request validation
package middlewares

import (
	"context"
	"net/http"

	"github.com/odundlaw/cbt-backend/internal/config"
	"github.com/odundlaw/cbt-backend/internal/constants"
	"github.com/odundlaw/cbt-backend/internal/helpers"
	"github.com/odundlaw/cbt-backend/internal/json"
	tokens "github.com/odundlaw/cbt-backend/internal/jwt"
	"github.com/odundlaw/cbt-backend/internal/store"
)

type contextKey string

const UserContextKey = contextKey("user")

func AuthMiddleware(next http.Handler, redis *store.Redis) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, _ := r.Cookie("access_token")
		if tokenStr.Value == "" {
			tokenStr.Value = helpers.BearerFromHeader(r)
		}

		if tokenStr.Value == "" {
			json.JSONError(w, http.StatusUnauthorized, constants.ErrUnauthorized, nil)
			return
		}

		claims, err := tokens.VerifyToken(tokenStr.Value, config.AccessSecret)
		if err != nil {
			json.JSONError(w, http.StatusUnauthorized, constants.ErrTokenInvalid, nil)
			return
		}

		if _, err := redis.GetJTI(r.Context(), claims.ID); err != nil {
			json.JSONError(w, http.StatusUnauthorized, constants.ErrTokenInvalid, nil)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
