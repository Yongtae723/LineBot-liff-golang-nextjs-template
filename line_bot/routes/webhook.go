package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"cookforyou.com/linebot-liff-template/line_bot/config"
	"cookforyou.com/linebot-liff-template/line_bot/logic/follow"
	"cookforyou.com/linebot-liff-template/line_bot/logic/message"

	"cookforyou.com/linebot-liff-template/common/llm"
	"cookforyou.com/linebot-liff-template/common/repository"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/rs/zerolog/log"
)

func HandleWebhook(c *gin.Context) {
	cfg := config.Load()
	bot, messageHandler, followHandler, err := setupDependencies(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	body, exists := c.Get("body")
	if !exists {
		log.Error().Msg("Request body not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	bodyBytes, ok := body.([]byte)
	if !ok {
		log.Error().Msg("Invalid body type")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var request webhook.CallbackRequest
	if err := json.Unmarshal(bodyBytes, &request); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal webhook request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	for _, event := range request.Events {
		switch e := event.(type) {
		case webhook.MessageEvent:
			handleMessageEvent(c.Request.Context(), e, bot, messageHandler)
		case webhook.FollowEvent:
			handleFollowEvent(c.Request.Context(), e, bot, followHandler)
		case webhook.UnfollowEvent:
			log.Info().Interface("event", e).Msg("Unfollow event received")
		default:
			log.Info().Interface("event", e).Msg("Unknown event type")
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func setupDependencies(cfg *config.Config) (*messaging_api.MessagingApiAPI, message.MessageHandler, follow.FollowHandler, error) {

	bot, err := messaging_api.NewMessagingApiAPI(cfg.LINE_CHANNEL_TOKEN)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create LINE bot client")
		return nil, nil, nil, err
	}

	userRepo := repository.NewUserRepo()
	convRepo := repository.NewConversationRepo()

	ctx := context.Background()
	geminiClient, err := llm.NewGoogleGemini(ctx, cfg.GEMINI_API_KEY, "gemini-2.5-flash-lite")
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Gemini client")
		return nil, nil, nil, err
	}

	messageHandler := message.NewMessageHandler(convRepo, userRepo, geminiClient)

	followHandler, err := follow.NewFollowHandler(cfg.LINE_CHANNEL_TOKEN, userRepo, cfg.LIFF_APP_URL, cfg.BACKEND_URL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create follow handler")
		return nil, nil, nil, err
	}

	return bot, messageHandler, followHandler, nil
}

func handleMessageEvent(ctx context.Context, event webhook.MessageEvent, bot *messaging_api.MessagingApiAPI, messageHandler message.MessageHandler) {
	switch message := event.Message.(type) {
	case webhook.TextMessageContent:
		lineID := event.Source.(webhook.UserSource).UserId
		responseText, err := messageHandler.HandleTextMessage(ctx, lineID, message.Text)
		if err != nil {
			log.Error().Err(err).Msg("Failed to handle text message")
			err := sendMessage(ctx, lineID, "エラーが発生しました。再度お試しください。", bot)
			if err != nil {
				log.Error().Err(err).Msg("Failed to send message")
			}
			return
		}
		if err := sendMessage(ctx, lineID, responseText, bot); err != nil {
			log.Error().Err(err).Msg("Failed to send message")
		}
	default:
		log.Info().Interface("message", message).Msg("Non-text message received")
	}
}

func handleFollowEvent(ctx context.Context, event webhook.FollowEvent, bot *messaging_api.MessagingApiAPI, followHandler follow.FollowHandler) {
	userID := event.Source.(webhook.UserSource).UserId
	welcomeMessage, err := followHandler.HandleFollow(ctx, userID)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("Failed to handle follow event")
		return
	}
	if err := sendMessage(ctx, userID, welcomeMessage, bot); err != nil {
		log.Error().Err(err).Msg("Failed to send welcome message")
	}
}

func sendMessage(ctx context.Context, userID, message string, bot *messaging_api.MessagingApiAPI) error {
	_, err := bot.PushMessage(&messaging_api.PushMessageRequest{
		To: userID,
		Messages: []messaging_api.MessageInterface{
			messaging_api.TextMessage{
				Text: message,
			},
		},
	}, "")

	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("Failed to send LINE message")
		return err
	}

	return nil
}
