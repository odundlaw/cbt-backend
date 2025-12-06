package users

import (
	"context"

	repo "github.com/odundlaw/cbt-backend/internal/adapters/postgresql/sqlc"
)

type Service interface {
	CreateUser(ctx context.Context, userParams createUserParams) (repo.User, error)
	CreateAdmin(ctx context.Context, adminParams createAdminParams) (repo.User, error)
	GetUserByID(ctx context.Context, ID int64) (repo.User, error)
	GetUserByEmail(ctx context.Context, email string) (repo.User, error)
	UpdateLastLogin(ctx context.Context, ID int64) (repo.User, error)
	UpdatePassword(ctx context.Context, params repo.UpdateUserPasswordParams) (repo.User, error)
}

type createUserParams struct {
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type loginParams struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"full_name" validate:"required,min=8"`
}

type forgotPasswordParams struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordParams struct {
	ResetToken string `json:"reset_token" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

type forgotPaswordResponse struct {
	Email             string `json:"email"`
	ResetTokenExpires string `json:"reset_token_expires"`
}

type UpdatePasswordResponse struct {
	UserID          int64  `json:"user_id"`
	PasswordResetAt string `json:"password_reset_at"`
}

type createAdminParams struct {
	FullName   string `json:"full_name" validate:"required,min=3,max=100"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	AdminCode  string `json:"admin_code" validate:"required"`
	Department string `json:"department" validate:"required"`
	Phone      string `json:"phone" validate:"required,min=11"`
}

type adminLoginParams struct {
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"full_name" validate:"required,min=8"`
	TwoFactorCode string `json:"two_factor_code"`
}

type adminForgotPasswordParams struct {
	Email     string `json:"email" validate:"required,email"`
	AdminCode string `json:"admin_code" validate:"required"`
}
