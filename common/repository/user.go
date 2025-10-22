package repository

import (
	"context"
	"errors"
	"fmt"

	"cookforyou.com/linebot-liff-template/common/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepo interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByLineID(ctx context.Context, lineID string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
}

type userRepo struct {
	*BaseRepo
}

func NewUserRepo() UserRepo {
	return &userRepo{BaseRepo: baseRepo}
}

func (r *userRepo) toMap(user *models.User) map[string]string {
	return map[string]string{
		"id":      user.ID,
		"line_id": user.LineID,
		"name":    user.Name,
	}
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	var users []*models.User
	_, err := r.Client.From("user").
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
	_, err := r.Client.From("user").
		Select("*", "", false).
		Eq("line_id", lineID).
		Limit(1, "").
		ExecuteTo(&users)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to get user by line_id: %w", err)
	}
	if len(users) == 0 {
		return nil, ErrUserNotFound
	}
	return users[0], nil
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	_, _, err := r.Client.From("user").
		Insert(r.toMap(user), false, "", "", "").
		Execute()
	if err != nil {
		return fmt.Errorf("repository: failed to create user: %w", err)
	}
	return nil
}
