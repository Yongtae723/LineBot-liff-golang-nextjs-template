package repository

import (
	"context"
	"fmt"

	"github.com/supabase-community/gotrue-go/types"
)

type AuthRepo interface {
	CreateUser(ctx context.Context, email string, password string, userMetadata map[string]interface{}, appMetadata map[string]interface{}) (string, error)
}

type authRepo struct {
	*BaseRepo
}

func NewAuthRepo() AuthRepo {
	return &authRepo{
		BaseRepo: baseRepo,
	}
}

func (r *authRepo) CreateUser(ctx context.Context, email string, password string, userMetadata map[string]interface{}, appMetadata map[string]interface{}) (string, error) {
	if _, ok := appMetadata["provider"]; !ok {
		appMetadata["provider"] = "line"
	}

	user := types.AdminCreateUserRequest{
		Email:        email,
		Password:     &password,
		UserMetadata: userMetadata,
		AppMetadata:  appMetadata,
		EmailConfirm: true,
	}

	response, err := r.Client.Auth.WithToken(r.BaseRepo.supabaseRoleKey).AdminCreateUser(user)
	if err != nil {
		return "", fmt.Errorf("auth_repo: failed to create user: %w", err)
	}

	return response.ID.String(), nil
}
