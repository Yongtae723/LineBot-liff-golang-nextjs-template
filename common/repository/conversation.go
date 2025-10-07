package repository

import (
	"context"
	"fmt"

	"cookforyou.com/linebot-liff-template/common/models"
	"github.com/supabase-community/postgrest-go"
)

type ConversationRepo interface {
	ListByUserID(ctx context.Context, userID string, limit int) ([]*models.Conversation, error)
	Create(ctx context.Context, conv *models.Conversation) error
}

type conversationRepo struct {
	*BaseRepo
}

func NewConversationRepo() ConversationRepo {
	return &conversationRepo{BaseRepo: baseRepo}
}

func (r *conversationRepo) toMap(conversation *models.Conversation) map[string]any {
	return map[string]any{
		"id":      conversation.ID,
		"user_id": conversation.UserID,
		"role":    conversation.Role,
		"content": conversation.Content,
	}
}

func (r *conversationRepo) ListByUserID(ctx context.Context, userID string, limit int) ([]*models.Conversation, error) {
	var conversations []*models.Conversation
	_, err := r.Client.From("conversations").
		Select("*", "", false).
		Eq("user_id", userID).
		Order("created_at", &postgrest.OrderOpts{Ascending: true}).
		Limit(limit, "").
		ExecuteTo(&conversations)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to get conversations: %w", err)
	}
	return conversations, nil
}

func (r *conversationRepo) Create(ctx context.Context, conv *models.Conversation) error {
	_, _, err := r.Client.From("conversations").
		Insert(r.toMap(conv), false, "", "", "").
		Execute()
	if err != nil {
		return fmt.Errorf("repository: failed to create conversation: %w", err)
	}
	return nil
}
