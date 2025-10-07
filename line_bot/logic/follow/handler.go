package follow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cookforyou.com/linebot-liff-template/line_bot/logic/follow/auth"

	"cookforyou.com/linebot-liff-template/common/repository"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/rs/zerolog/log"
)

type FollowHandler interface {
	HandleFollow(ctx context.Context, userID string) (string, error)
}

type followHandler struct {
	bot        *messaging_api.MessagingApiAPI
	userRepo   repository.UserRepo
	liffAppURL string
	backendURL string
}

func NewFollowHandler(channelToken string, userRepo repository.UserRepo, liffAppURL, backendURL string) (FollowHandler, error) {
	bot, err := messaging_api.NewMessagingApiAPI(channelToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create LINE bot client: %w", err)
	}

	return &followHandler{
		bot:        bot,
		userRepo:   userRepo,
		liffAppURL: liffAppURL,
		backendURL: backendURL,
	}, nil
}

func (h *followHandler) HandleFollow(ctx context.Context, userID string) (string, error) {
	user, err := h.userRepo.GetByLineID(ctx, userID)
	if err != nil && err != repository.ErrUserNotFound {
		return "", fmt.Errorf("follow_handler: failed to get user: %w", err)
	}

	if user != nil {
		return h.buildWelcomeMessage(user.Name), nil
	}

	profile, err := h.bot.GetProfile(userID)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("follow_handler: Failed to get user profile")
		return "", fmt.Errorf("follow_handler: failed to get user profile: %w", err)
	}
	// Register user in backend
	if err := h.registerUserInBackend(ctx, userID, profile.DisplayName); err != nil {
		return "", fmt.Errorf("follow_handler: failed to register user in backend: %w", err)
	}

	return h.buildWelcomeMessage(profile.DisplayName), nil
}

func (h *followHandler) buildWelcomeMessage(displayName string) string {
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

type BotUserRegisterRequest struct {
	LineID      string `json:"line_id"`
	DisplayName string `json:"display_name"`
}

func (h *followHandler) registerUserInBackend(ctx context.Context, lineID, displayName string) error {
	token := auth.GenerateJWTToken(lineID)
	reqBody := BotUserRegisterRequest{
		LineID:      lineID,
		DisplayName: displayName,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/user/register/bot", h.backendURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("backend returned status %d", resp.StatusCode)
	}

	log.Info().
		Str("user_id", lineID).
		Str("display_name", displayName).
		Msg("Successfully registered user in backend")

	return nil
}
