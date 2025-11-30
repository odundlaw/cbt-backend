// Package users where all users related information is handled
package users

import (
	"context"
	"errors"

	repo "github.com/odundlaw/cbt-backend/internal/adapters/postgresql/sqlc"
	"github.com/odundlaw/cbt-backend/internal/constants"
	"github.com/odundlaw/cbt-backend/internal/helpers"
)

type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries) Service {
	return &svc{repo: repo}
}

func (s *svc) CreateUser(ctx context.Context, userParams createUserParams) (repo.User, error) {
	hashed, err := helpers.HashPassword(userParams.Password)
	if err != nil {
		return repo.User{}, errors.New(constants.ErrFailedHashPass)
	}

	user := repo.CreateUserParams{
		FullName: userParams.FullName,
		Email:    userParams.Email,
		Password: hashed,
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *svc) CreateAdmin() {}

func (s *svc) GetUserByID(ctx context.Context, ID int64) (repo.User, error) {
	return s.repo.GetUserByID(ctx, ID)
}

func (s *svc) GetUserByEmail(ctx context.Context, email string) (repo.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *svc) UpdateLastLogin(ctx context.Context, ID int64) (repo.User, error) {
	return s.repo.UpdateLastLogin(ctx, ID)
}

func (s *svc) UpdatePassword(ctx context.Context, params repo.UpdateUserPasswordParams) (repo.User, error) {
	hashed, err := helpers.HashPassword(params.Password)
	if err != nil {
		return repo.User{}, errors.New(constants.ErrFailedHashPass)
	}

	update := repo.UpdateUserPasswordParams{
		ID:       params.ID,
		Password: hashed,
	}

	return s.repo.UpdateUserPassword(ctx, update)
}
