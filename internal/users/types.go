package users

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	repo "github.com/odundlaw/cbt-backend/internal/adapters/postgresql/sqlc"
	response "github.com/odundlaw/cbt-backend/internal/json"
)

type Service interface {
	CreateUser(ctx context.Context, userParams createUserParams) (repo.User, error)
	CreateAdmin()
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
	ResetTokenExpires int    `json:"reset_token_expires"`
}

type UpdatePasswordResponse struct {
	UserID          int64  `json:"user_id"`
	PasswordResetAt string `json:"password_reset_at"`
}

type createAdminParams struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func FormatValidationErrors(err error) []response.FieldError {
	errs := []response.FieldError{}

	for _, fe := range err.(validator.ValidationErrors) {
		msg := ""

		switch fe.Tag() {
		case "required":
			msg = fmt.Sprintf("%s is required", fe.Field())
		case "email":
			msg = "Invalid email format"
		case "min":
			msg = fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
		case "max":
			msg = fmt.Sprintf("%s must not be more than %s characters", fe.Field(), fe.Param())
		default:
			msg = fmt.Sprintf("%s is not valid", fe.Field())
		}

		errs = append(errs, response.FieldError{
			Field:   fe.Field(),
			Message: msg,
		})
	}

	return errs
}
