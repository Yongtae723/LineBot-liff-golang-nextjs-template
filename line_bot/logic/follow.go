package logic

import (
	"context"
	"fmt"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/rs/zerolog/log"
)

type FollowHandler struct {
	bot        *messaging_api.MessagingApiAPI
	liffAppURL string
}

func NewFollowHandler(channelToken, liffAppURL string) (*FollowHandler, error) {
	bot, err := messaging_api.NewMessagingApiAPI(channelToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create LINE bot client: %w", err)
	}

	return &FollowHandler{
		bot:        bot,
		liffAppURL: liffAppURL,
	}, nil
}

func (h *FollowHandler) HandleFollow(ctx context.Context, userID string) (string, error) {
	profile, err := h.bot.GetProfile(userID)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("Failed to get user profile")
		return h.buildWelcomeMessage(""), nil
	}

	log.Info().
		Str("user_id", userID).
		Str("display_name", profile.DisplayName).
		Msg("New follower")

	return h.buildWelcomeMessage(profile.DisplayName), nil
}

func (h *FollowHandler) buildWelcomeMessage(displayName string) string {
	if displayName != "" {
		return fmt.Sprintf(
			"こんにちは、%sさん！\n友だち追加ありがとうございます🎉\n\n"+
				"このBotはGemini AIと会話できるチャットボットです。\n\n"+
				"📱 初回登録が必要です：\n%s\n\n"+
				"登録後、メッセージを送ってみてください！",
			displayName,
			h.liffAppURL,
		)
	}
	return fmt.Sprintf(
		"友だち追加ありがとうございます🎉\n\n"+
			"このBotはGemini AIと会話できるチャットボットです。\n\n"+
			"📱 初回登録が必要です：\n%s\n\n"+
			"登録後、メッセージを送ってみてください！",
		h.liffAppURL,
	)
}
