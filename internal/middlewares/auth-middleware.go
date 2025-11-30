// Package middlewares where all middlware files are stored for proper request validation
package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/odundlaw/cbt-backend/internal/config"
	"github.com/odundlaw/cbt-backend/internal/constants"
	"github.com/odundlaw/cbt-backend/internal/json"
	tokens "github.com/odundlaw/cbt-backend/internal/jwt"
)

type contextKey string

const UserContextKey = contextKey("user")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorizaion")
		refreshHeader := r.Header.Get("X-Refresh-Token")
		if authHeader == "" {
			json.JSONError(w, http.StatusUnauthorized, constants.ErrUnauthorized, nil)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			json.JSONError(w, http.StatusUnauthorized, constants.ErrInvalidAuthHeader, nil)
			return
		}

		accessToken := parts[1]

		claims, err := tokens.VerifyToken(accessToken, config.AccessSecret)
		if err == nil {
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if errors.Is(err, jwt.ErrTokenExpired) {
			if refreshHeader == "" {
				json.JSONError(w, http.StatusUnauthorized, constants.ErrInvalidRefreshToken, nil)
				return
			}

			refreshClaims, err := tokens.VerifyToken(refreshHeader, config.RefreshSecret)
			if err != nil {
				json.JSONError(w, http.StatusUnauthorized, constants.ErrInvalidRefreshToken, nil)
				return
			}

			newTokens, err := tokens.GenerateTokens(refreshClaims.UserID, refreshClaims.Email)
			if err != nil {
				json.JSONError(w, http.StatusUnauthorized, constants.ErrFailedTokenGen, nil)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, refreshClaims)

			w.Header().Set("X-New-Access-Token", newTokens.AccessToken)
			w.Header().Set("X-New-Refresh-Token", newTokens.RefreshToken)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		json.JSONError(w, http.StatusUnauthorized, constants.ErrTokenInvalid, nil)
	})
}
