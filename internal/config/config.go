// Package config, where all environment config are stored
package config

import "github.com/odundlaw/cbt-backend/internal/env"

var (
	AccessSecret        = []byte(env.GetString("ACCESS_TOKEN_SECRET", ""))
	RefreshSecret       = []byte(env.GetString("REFRESH_TOKEN_SECRET", ""))
	ResetPasswordSecret = []byte(env.GetString("RESET_PASSWORD_SECRET", ""))
	DatabaseURL         = env.GetString("DATABASE_URL", "")
)
