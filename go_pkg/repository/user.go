package repository

import (
	"context"
	"fmt"

	"github.com/linebot-liff-template/go_pkg/models"
)

type UserRepo interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByLineID(ctx context.Context, lineID string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
}

type userRepo struct {
	*BaseRepo
}

func NewUserRepo() UserRepo {
	return &userRepo{BaseRepo: baseRepo}
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	var users []*models.User
	_, err := r.Client.From("users").
		Select("*", "", false).
		Eq("id", id).
		Limit(1, "").
		ExecuteTo(&users)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to get user by id: %w", err)
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("repository: user not found")
	}
	return users[0], nil
}

func (r *userRepo) GetByLineID(ctx context.Context, lineID string) (*models.User, error) {
	var users []*models.User
	_, err := r.Client.From("users").
		Select("*", "", false).
		Eq("line_id", lineID).
		Limit(1, "").
		ExecuteTo(&users)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to get user by line_id: %w", err)
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("repository: user not found")
	}
	return users[0], nil
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	_, _, err := r.Client.From("users").
		Insert(user, false, "", "", "").
		Execute()
	if err != nil {
		return fmt.Errorf("repository: failed to create user: %w", err)
	}
	return nil
}

func (r *userRepo) Update(ctx context.Context, user *models.User) error {
	_, _, err := r.Client.From("users").
		Update(user, "", "").
		Eq("id", user.ID).
		Execute()
	if err != nil {
		return fmt.Errorf("repository: failed to update user: %w", err)
	}
	return nil
}
