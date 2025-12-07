package users

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	repo "github.com/odundlaw/cbt-backend/internal/adapters/postgresql/sqlc"
	"github.com/odundlaw/cbt-backend/internal/config"
	"github.com/odundlaw/cbt-backend/internal/constants"
	"github.com/odundlaw/cbt-backend/internal/helpers"
	"github.com/odundlaw/cbt-backend/internal/json"
	"github.com/odundlaw/cbt-backend/internal/jwt"
	"github.com/odundlaw/cbt-backend/internal/store"
	"github.com/odundlaw/cbt-backend/internal/validation"
)

type Handler struct {
	service Service
	rdb     *store.Redis
}

func NewHandler(service Service, rdb *store.Redis) *Handler {
	return &Handler{
		service,
		rdb,
	}
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
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

	if err := jwt.Persist(r.Context(), h.rdb, tokens); err != nil {
		json.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helpers.SetAuthCookies(w, tokens)

	json.JSONSuccess(w, http.StatusOK, constants.MsgAccountCreated, user, &json.Token{
		AccessToken: tokens.Access,
		ExpiresIn:   tokens.ExpAcc.Second(),
	})
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
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

	if err := jwt.Persist(r.Context(), h.rdb, tokens); err != nil {
		json.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helpers.SetAuthCookies(w, tokens)

	json.JSONSuccess(w, http.StatusOK, constants.MsgLoginSuccessful, user, &json.Token{
		AccessToken: tokens.Access,
		ExpiresIn:   tokens.ExpAcc.Second(),
	})
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
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

	err = jwt.PersistResetToken(r.Context(), h.rdb, token.JTIRes, strconv.FormatInt(user.ID, 10), token.ExpRes)
	if err != nil {
		json.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// todo send Email Message

	res := forgotPaswordResponse{
		Email:             user.Email,
		ResetTokenExpires: token.ExpRes.Local().String(),
	}

	json.JSONSuccess(w, http.StatusOK, constants.PswResetSentSuccessful, res, nil)
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) LoginAdmin(w http.ResponseWriter, r *http.Request) {
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

	helpers.SetAuthCookies(w, tokens)

	json.JSONSuccess(w, http.StatusOK, constants.MsgAdminLoginSuccessful, admin, &json.Token{
		AccessToken: tokens.Access,
		ExpiresIn:   tokens.ExpAcc.Second(),
	})
}

func (h *Handler) AdminForgotPassword(w http.ResponseWriter, r *http.Request) {
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
		ResetTokenExpires: token.ExpRes.Local().String(),
	}

	json.JSONSuccess(w, http.StatusOK, constants.PswResetSentSuccessful, res, nil)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	acc, _ := r.Cookie("access_token")
	ref, _ := r.Cookie("refresh_token")

	var userID int64

	if acc.Value != "" {
		if claims, err := jwt.VerifyToken(acc.Value, config.AccessSecret); err != nil {
			userID = claims.UserID
			_ = h.rdb.DelJTI(r.Context(), "access:"+claims.ID)
		}
	}

	if ref.Value != "" {
		if claims, err := jwt.VerifyToken(ref.Value, config.RefreshSecret); err != nil {
			_ = h.rdb.DelJTI(r.Context(), "refresh:"+claims.ID)
		}
	}

	helpers.ClearAuthCookies(w)
	json.JSONSuccess(w, http.StatusOK, constants.MsgLogoutSuccessful, LogoutResData{
		UserID:      userID,
		LoggedOutAt: time.Now().String(),
	}, nil)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ref, err := jwt.MustCookie(r, "refresh_token")
	if err != nil {
		json.JSONError(w, http.StatusUnauthorized, constants.ErrRefreshTokenRequired, nil)
		return
	}

	claims, err := jwt.VerifyToken(ref, config.RefreshSecret)
	if err != nil {
		json.JSONError(w, http.StatusUnauthorized, constants.ErrInvalidRefreshToken, nil)
		return
	}

	userID := strconv.FormatInt(claims.UserID, 10)

	if _, err := h.rdb.GetJTI(r.Context(), "refresh:"+userID); err != nil {
		json.JSONError(w, http.StatusUnauthorized, constants.ErrRevokedToken, nil)
		return
	}
	_ = h.rdb.DelJTI(r.Context(), "refresh:"+userID)

	toks, err := jwt.GenerateTokens(claims.UserID, claims.Email)
	if err != nil {
		json.JSONError(w, http.StatusInternalServerError, constants.ErrFailedTokenGen, nil)
		return
	}

	if err := jwt.Persist(r.Context(), h.rdb, toks); err != nil {
		json.JSONError(w, http.StatusInternalServerError, constants.ErrPersistToken, nil)
		return
	}

	helpers.SetAuthCookies(w, toks)
	json.JSONSuccess(w, http.StatusOK, constants.MsgRefresSuccessful, nil, &json.Token{
		AccessToken: toks.Access,
		ExpiresIn:   toks.ExpAcc.Second(),
	})
}
