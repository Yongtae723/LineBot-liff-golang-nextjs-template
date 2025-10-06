package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/linebot-liff-template/line_bot/logic"
	"github.com/rs/zerolog/log"
)

type WebhookHandler struct {
	bot            *messaging_api.MessagingApiAPI
	messageHandler *logic.MessageHandler
	followHandler  *logic.FollowHandler
}

func NewWebhookHandler(bot *messaging_api.MessagingApiAPI, messageHandler *logic.MessageHandler, followHandler *logic.FollowHandler) *WebhookHandler {
	return &WebhookHandler{
		bot:            bot,
		messageHandler: messageHandler,
		followHandler:  followHandler,
	}
}

func (h *WebhookHandler) Handle(c *gin.Context) {
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
			h.handleMessageEvent(c.Request.Context(), e)
		case webhook.FollowEvent:
			h.handleFollowEvent(c.Request.Context(), e)
		case webhook.UnfollowEvent:
			log.Info().Interface("event", e).Msg("Unfollow event received")
		default:
			log.Info().Interface("event", e).Msg("Unknown event type")
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *WebhookHandler) handleMessageEvent(ctx context.Context, event webhook.MessageEvent) {
	switch message := event.Message.(type) {
	case webhook.TextMessageContent:
		userID := event.Source.(webhook.UserSource).UserId
		responseText, err := h.messageHandler.HandleTextMessage(ctx, userID, message.Text)
		if err != nil {
			log.Error().Err(err).Msg("Failed to handle text message")
			return
		}
		if err := h.sendMessage(ctx, userID, responseText); err != nil {
			log.Error().Err(err).Msg("Failed to send message")
		}
	default:
		log.Info().Interface("message", message).Msg("Non-text message received")
	}
}

func (h *WebhookHandler) handleFollowEvent(ctx context.Context, event webhook.FollowEvent) {
	userID := event.Source.(webhook.UserSource).UserId
	welcomeMessage, err := h.followHandler.HandleFollow(ctx, userID)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("Failed to handle follow event")
		return
	}
	if err := h.sendMessage(ctx, userID, welcomeMessage); err != nil {
		log.Error().Err(err).Msg("Failed to send welcome message")
	}
}

func (h *WebhookHandler) sendMessage(ctx context.Context, userID, message string) error {
	_, err := h.bot.PushMessage(&messaging_api.PushMessageRequest{
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
