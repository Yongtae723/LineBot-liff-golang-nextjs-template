package handler

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/linebot-liff-template/go_pkg/llm"
	"github.com/linebot-liff-template/go_pkg/models"
	"github.com/linebot-liff-template/go_pkg/repository"
	"github.com/rs/zerolog/log"
)

type MessageHandler struct {
	bot          *messaging_api.MessagingApiAPI
	convRepo     repository.ConversationRepo
	userRepo     repository.UserRepo
	geminiClient llm.GoogleGemini
}

func NewMessageHandler(channelToken, geminiAPIKey string) (*MessageHandler, error) {
	bot, err := messaging_api.NewMessagingApiAPI(channelToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create LINE bot client: %w", err)
	}

	ctx := context.Background()
	geminiClient, err := llm.NewGoogleGemini(ctx, geminiAPIKey, "gemini-1.5-flash")
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &MessageHandler{
		bot:          bot,
		convRepo:     repository.NewConversationRepo(),
		userRepo:     repository.NewUserRepo(),
		geminiClient: geminiClient,
	}, nil
}

func (h *MessageHandler) HandleTextMessage(ctx context.Context, userID, messageText string) error {
	user, err := h.userRepo.GetByLineID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("User not found")
		return h.replyMessage(ctx, userID, "ユーザー登録が必要です。LIFFアプリから登録してください。")
	}

	userConv := &models.Conversation{
		ID:      uuid.New().String(),
		UserID:  user.LineID,
		Role:    models.RoleUser,
		Content: messageText,
	}
	if err := h.convRepo.Create(ctx, userConv); err != nil {
		log.Error().Err(err).Msg("Failed to save user message")
		return err
	}

	history, err := h.convRepo.ListByUserID(ctx, user.LineID, 20)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get conversation history")
		return err
	}

	response, err := h.geminiClient.Chat(ctx, history)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get LLM response")
		return h.replyMessage(ctx, userID, "申し訳ございません。エラーが発生しました。")
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

	return h.replyMessage(ctx, userID, response)
}

func (h *MessageHandler) replyMessage(ctx context.Context, userID, message string) error {
	_, err := h.bot.PushMessage(&messaging_api.PushMessageRequest{
		To: userID,
		Messages: []messaging_api.MessageInterface{
			messaging_api.TextMessage{
				Text: message,
			},
		},
	}, "")

	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("Failed to send message")
		return err
	}

	return nil
}
