package message

import (
	"context"

	"cookforyou.com/linebot-liff-template/common/llm"
	"cookforyou.com/linebot-liff-template/common/models"
	"cookforyou.com/linebot-liff-template/common/repository"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type MessageHandler interface {
	HandleTextMessage(ctx context.Context, userID, messageText string) (string, error)
}

type messageHandler struct {
	convRepo     repository.ConversationRepo
	userRepo     repository.UserRepo
	geminiClient llm.GoogleGemini
}

func NewMessageHandler(convRepo repository.ConversationRepo, userRepo repository.UserRepo, geminiClient llm.GoogleGemini) MessageHandler {
	return &messageHandler{
		convRepo:     convRepo,
		userRepo:     userRepo,
		geminiClient: geminiClient,
	}
}

func (h *messageHandler) HandleTextMessage(ctx context.Context, userID, messageText string) (string, error) {
	user, err := h.userRepo.GetByLineID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("User not found")
		return "ユーザー登録が必要です。LIFFアプリから登録してください。", nil
	}

	userConv := &models.Conversation{
		ID:      uuid.New().String(),
		UserID:  user.LineID,
		Role:    models.RoleUser,
		Content: messageText,
	}
	if err := h.convRepo.Create(ctx, userConv); err != nil {
		log.Error().Err(err).Msg("Failed to save user message")
		return "", err
	}

	history, err := h.convRepo.ListByUserID(ctx, user.LineID, 20)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get conversation history")
		return "", err
	}

	response, err := h.geminiClient.Chat(ctx, history)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get LLM response")
		return "申し訳ございません。エラーが発生しました。", nil
	}

	assistantConv := &models.Conversation{
		ID:      uuid.New().String(),
		UserID:  user.LineID,
		Role:    models.RoleAssistant,
		Content: response,
	}
	if err := h.convRepo.Create(ctx, assistantConv); err != nil {
		log.Error().Err(err).Msg("Failed to save assistant message")
	}

	return response, nil
}
