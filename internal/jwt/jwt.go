// Package jwt for managing generating and invalidating jwt tokens
package jwt

import (
	"context"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/odundlaw/cbt-backend/internal/config"
	"github.com/odundlaw/cbt-backend/internal/store"
)

type Claims struct {
	UserID int64
	Email  string
	jwt.RegisteredClaims
}

type Tokens struct {
	Access     string
	Refresh    string
	ResetToken string
	JTIAcc     string
	JTIRef     string
	JTIRes     string
	ExpAcc     time.Time
	ExpRef     time.Time
	ExpRes     time.Time
	UserID     int64
	Email      string
	Issuer     string
	Audience   string
}

func GenerateTokens(userID int64, email string) (*Tokens, error) {
	now := time.Now().UTC()

	t := &Tokens{
		UserID:   userID,
		Email:    email,
		JTIAcc:   uuid.NewString(),
		JTIRef:   uuid.NewString(),
		ExpAcc:   now.Add(15 * time.Minute),
		ExpRef:   now.Add(7 * 24 * time.Hour),
		Issuer:   "cbt-backend",
		Audience: "cbt-users",
	}

	accessClaims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t.ExpAcc),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        t.JTIAcc,
			Issuer:    t.Issuer,
			Subject:   strconv.FormatInt(int64(userID), 10),
			Audience:  jwt.ClaimStrings{t.Audience},
		},
	}

	refreshClaim := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t.ExpRef),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        t.JTIRef,
			Issuer:    t.Issuer,
			Subject:   strconv.FormatInt(int64(userID), 10),
			Audience:  jwt.ClaimStrings{t.Audience},
		},
	}

	var err error
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	t.Access, err = accessToken.SignedString(config.AccessSecret)
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	t.Refresh, err = refreshToken.SignedString(config.RefreshSecret)
	if err != nil {
		return nil, err
	}

	return t, nil
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

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, jwt.ErrTokenExpired
	}

	return claims, nil
}

func Persist(ctx context.Context, r *store.Redis, t *Tokens) error {
	if err := r.SetJTI(ctx, "access:"+t.JTIAcc, strconv.FormatInt(int64(t.UserID), 10), t.ExpAcc); err != nil {
		return err
	}

	if err := r.SetJTI(ctx, "refresh:"+t.JTIRef, strconv.FormatInt(int64(t.UserID), 10), t.ExpRef); err != nil {
		return err
	}

	return nil
}

func PersistResetToken(ctx context.Context, r *store.Redis, jti, userID string, exp time.Time) error {
	if err := r.SetJTI(ctx, "reset:"+jti, userID, exp); err != nil {
		return err
	}
	return nil
}

func GenerateResetPasswordToken(userID int64, email string) (*Tokens, error) {
	now := time.Now().UTC()

	t := &Tokens{
		UserID:   userID,
		Email:    email,
		JTIRes:   uuid.NewString(),
		ExpRes:   now.Add(15 * time.Minute),
		Issuer:   "cbt-backend",
		Audience: "cbt-users",
	}

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t.ExpRes),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   strconv.FormatInt(int64(userID), 10),
			ID:        t.JTIRes,
			Issuer:    t.Issuer,
			Audience:  jwt.ClaimStrings{t.Audience},
		},
	}

	var err error
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t.ResetToken, err = token.SignedString(config.ResetPasswordSecret)
	if err != nil {
		return nil, err
	}

	return t, nil
}
