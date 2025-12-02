package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	repo "github.com/odundlaw/cbt-backend/internal/adapters/postgresql/sqlc"
	"github.com/odundlaw/cbt-backend/internal/config"
	"github.com/odundlaw/cbt-backend/internal/constants"
	"github.com/odundlaw/cbt-backend/internal/helpers"
	"github.com/odundlaw/cbt-backend/internal/json"
	"github.com/odundlaw/cbt-backend/internal/jwt"
	"github.com/odundlaw/cbt-backend/internal/validation"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service,
	}
}

func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req createUserParams

	if err := json.ReadJSON(r, &req); err != nil {
		json.JSONError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		formattedErr := validation.FormatValidationErrors(err)
		json.JSONError(w, http.StatusBadRequest, constants.ErrValidationFailed, formattedErr)
		return
	}

	user, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		json.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	tokens, err := jwt.GenerateTokens(user.ID, user.Email)
	if err != nil {
		json.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	go func(userID int64) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, err := h.service.UpdateLastLogin(ctx, userID)
		if err != nil {
			fmt.Println("failed to update last login:", err)
		}
	}(user.ID)

	json.JSONSuccess(w, http.StatusOK, constants.MsgAccountCreated, user, tokens)
}

func (h *handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req loginParams

	if err := json.ReadJSON(r, &req); err != nil {
		json.JSONError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		formattedErr := validation.FormatValidationErrors(err)
		json.JSONError(w, http.StatusBadRequest, constants.ErrValidationFailed, formattedErr)
		return
	}

	user, err := h.service.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		json.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if err := helpers.CheckPassword(user.Password, req.Password); err != nil {
		json.JSONError(w, http.StatusBadRequest, constants.ErrInvalidLogin, nil)
		return
	}

	tokens, err := jwt.GenerateTokens(user.ID, user.Email)
	if err != nil {
		json.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	go func(userID int64) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, err := h.service.UpdateLastLogin(ctx, userID)
		if err != nil {
			fmt.Println("failed to update last login:", err)
		}
	}(user.ID)

	json.JSONSuccess(w, http.StatusOK, constants.MsgLoginSuccessful, user, tokens)
}

func (h *handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req forgotPasswordParams

	if err := json.ReadJSON(r, &req); err != nil {
		json.JSONError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		formattedErr := validation.FormatValidationErrors(err)
		json.JSONError(w, http.StatusBadRequest, constants.ErrValidationFailed, formattedErr)
		return
	}

	user, err := h.service.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		json.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	token, err := jwt.GenerateResetPasswordToken(user.ID, user.Email)
	if err != nil {
		json.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// todo send Email Message

	res := forgotPaswordResponse{
		Email:             user.Email,
		ResetTokenExpires: token.Expires,
	}

	json.JSONSuccess(w, http.StatusOK, constants.PswResetSentSuccessful, res, nil)
}

func (h *handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordParams

	if err := json.ReadJSON(r, &req); err != nil {
		json.JSONError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		formattedErr := validation.FormatValidationErrors(err)
		json.JSONError(w, http.StatusBadRequest, constants.ErrValidationFailed, formattedErr)
		return
	}

	claims, err := jwt.VerifyToken(req.ResetToken, config.ResetPasswordSecret)
	if err != nil {
		json.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	user, err := h.service.UpdatePassword(r.Context(), repo.UpdateUserPasswordParams{
		ID:       claims.UserID,
		Password: req.Password,
	})
	if err != nil {
		json.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	json.JSONSuccess(w, http.StatusOK, constants.PswResetSentSuccessful, UpdatePasswordResponse{
		UserID:          user.ID,
		PasswordResetAt: user.UpdatedAt.Time.String(),
	}, nil)
}

func (h *handler) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	var req createAdminParams

	if err := json.ReadJSON(r, &req); err != nil {
		json.JSONError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		formattedErr := validation.FormatValidationErrors(err)
		json.JSONError(w, http.StatusBadRequest, constants.ErrValidationFailed, formattedErr)
		return
	}

	// Note to work on validating Admin_code
	// Remember to add Admin Permissions

	adminUser, err := h.service.CreateAdmin(r.Context(), req)
	if err != nil {
		json.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	json.JSONSuccess(w, http.StatusOK, constants.MsgAdminAccountCreated, adminUser, nil)
}

func (h *handler) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	var req adminLoginParams

	if err := json.ReadJSON(r, &req); err != nil {
		json.JSONError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		formattedErr := validation.FormatValidationErrors(err)
		json.JSONError(w, http.StatusBadRequest, constants.ErrValidationFailed, formattedErr)
		return
	}

	admin, err := h.service.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		json.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if err := helpers.CheckPassword(admin.Password, req.Password); err != nil {
		json.JSONError(w, http.StatusBadRequest, constants.ErrInvalidLogin, nil)
		return
	}

	if admin.Status != repo.UserStatusApproved {
		json.JSONError(w, http.StatusBadRequest, constants.ErrAccountNotApporve, nil)
		return
	}

	tokens, err := jwt.GenerateTokens(admin.ID, admin.Email)
	if err != nil {
		json.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	go func(userID int64) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, err := h.service.UpdateLastLogin(ctx, userID)
		if err != nil {
			fmt.Println("failed to update last login:", err)
		}
	}(admin.ID)

	json.JSONSuccess(w, http.StatusOK, constants.MsgAdminLoginSuccessful, admin, tokens)
}

func (h *handler) AdminForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req adminForgotPasswordParams

	if err := json.ReadJSON(r, &req); err != nil {
		json.JSONError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		formattedErr := validation.FormatValidationErrors(err)
		json.JSONError(w, http.StatusBadRequest, constants.ErrValidationFailed, formattedErr)
		return
	}

	user, err := h.service.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		json.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	token, err := jwt.GenerateResetPasswordToken(user.ID, user.Email)
	if err != nil {
		json.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// todo send Email Message

	res := forgotPaswordResponse{
		Email:             user.Email,
		ResetTokenExpires: token.Expires,
	}

	json.JSONSuccess(w, http.StatusOK, constants.PswResetSentSuccessful, res, nil)
}
