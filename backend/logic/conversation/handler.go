package conversation

import (
	"context"
	"fmt"

	"cookforyou.com/linebot-liff-template/common/llm"
	"cookforyou.com/linebot-liff-template/common/models"
	"cookforyou.com/linebot-liff-template/common/repository"
	"github.com/google/uuid"
)

type Handler struct {
	convRepo  repository.ConversationRepo
	userRepo  repository.UserRepo
	llmClient llm.GoogleGemini
}

func NewHandler(convRepo repository.ConversationRepo, userRepo repository.UserRepo, llmClient llm.GoogleGemini) *Handler {
	return &Handler{
		convRepo:  convRepo,
		userRepo:  userRepo,
		llmClient: llmClient,
	}
}

func (h *Handler) GetHistory(ctx context.Context, userID string, limit int) ([]*models.Conversation, error) {
	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	conversations, err := h.convRepo.ListByUserID(ctx, user.LineID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversations: %w", err)
	}

	return conversations, nil
}

func (h *Handler) ProcessMessage(ctx context.Context, userID string, message string) (string, error) {
	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	userConv := &models.Conversation{
		ID:      uuid.New().String(),
		UserID:  user.LineID,
		Role:    models.RoleUser,
		Content: message,
	}
	if err := h.convRepo.Create(ctx, userConv); err != nil {
		return "", fmt.Errorf("failed to save user message: %w", err)
	}

	history, err := h.convRepo.ListByUserID(ctx, user.LineID, 20)
	if err != nil {
		return "", fmt.Errorf("failed to get history: %w", err)
	}

	response, err := h.llmClient.Chat(ctx, history)
	if err != nil {
		return "", fmt.Errorf("failed to get LLM response: %w", err)
	}

	assistantConv := &models.Conversation{
		ID:      uuid.New().String(),
		UserID:  user.LineID,
		Role:    models.RoleAssistant,
		Content: response,
	}
	if err := h.convRepo.Create(ctx, assistantConv); err != nil {
		return "", fmt.Errorf("failed to save assistant message: %w", err)
	}

	return response, nil
}
