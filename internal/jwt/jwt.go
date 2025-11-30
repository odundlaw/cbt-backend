// Package jwt for managing generating and invalidating jwt tokens
package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/odundlaw/cbt-backend/internal/config"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type ForgotPasswordToken struct {
	ResetToken string `json:"reset_token"`
	Expires    int    `json:"reset_token_expires"`
}

func GenerateTokens(userID int64, email string) (*Token, error) {
	accessExp := time.Now().Add(1 * time.Hour)

	accessClaims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccessToken, err := accessToken.SignedString(config.AccessSecret)
	if err != nil {
		return nil, err
	}

	refreshExp := time.Now().Add(7 * 24 * time.Hour)

	refreshClaim := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	signedRefreshToken, err := refreshToken.SignedString(config.RefreshSecret)
	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
		ExpiresIn:    int(time.Until(accessExp).Seconds()),
	}, nil
}

func VerifyToken(tokenStr string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidId
	}

	return claims, nil
}

func GenerateResetPasswordToken(userID int64, email string) (*ForgotPasswordToken, error) {
	exp := time.Now().Add(15 * time.Minute)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(config.ResetPasswordSecret)
	if err != nil {
		return nil, err
	}

	return &ForgotPasswordToken{
		ResetToken: signedToken,
		Expires:    int(time.Until(exp).Seconds()),
	}, nil
}
