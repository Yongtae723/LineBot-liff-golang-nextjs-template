package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cookforyou.com/linebot-liff-template/backend/config"
	"cookforyou.com/linebot-liff-template/common/models"
	"cookforyou.com/linebot-liff-template/common/repository"
)

type RegisterMethod interface {
	RegisterFromAccessToken(ctx context.Context, accessToken string) (string, error)
	Register(ctx context.Context, lineID, displayName string) (string, error)
}

type registerHandler struct {
	userRepo      repository.UserRepo
	authRepo      repository.AuthRepo
	lineChannelID string
}

func NewRegisterHandler(userRepo repository.UserRepo, authRepo repository.AuthRepo, lineChannelID string) RegisterMethod {
	return &registerHandler{
		userRepo:      userRepo,
		authRepo:      authRepo,
		lineChannelID: lineChannelID,
	}
}

type LineProfile struct {
	UserID      string `json:"userId"`
	DisplayName string `json:"displayName"`
	PictureURL  string `json:"pictureUrl"`
}

func (h *registerHandler) RegisterFromAccessToken(ctx context.Context, accessToken string) (string, error) {
	lineProfile, err := h.getLineProfile(accessToken)
	if err != nil {
		return "", fmt.Errorf("failed to get LINE profile: %w", err)
	}

	return h.Register(ctx, lineProfile.UserID, lineProfile.DisplayName)
}

func (h *registerHandler) Register(ctx context.Context, lineID, displayName string) (string, error) {

	// NOTE: This is a simplified authentication flow for template purposes.
	// In production, implement proper authentication with:
	// - LINE Login integration with Supabase Auth
	// - Proper OAuth flow
	// - Secure token exchange
	// - Email verification (if needed)
	email := lineID + "@example.com"
	password := lineID

	userMetadata := map[string]interface{}{
		"line_id": lineID,
		"name":    displayName,
	}
	appMetadata := map[string]interface{}{
		"provider": "line",
	}

	authUserID, err := h.authRepo.CreateUser(ctx, email, password, userMetadata, appMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to create auth user: %w", err)
	}

	newUser := &models.User{
		ID:     authUserID,
		LineID: lineID,
		Name:   displayName,
	}

	if err := h.userRepo.Create(ctx, newUser); err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return lineID, nil
}

func (h *registerHandler) getLineProfile(accessToken string) (*LineProfile, error) {
	if accessToken == "local_access_token" {
		cfg := config.Load()
		return &LineProfile{
			UserID:      cfg.MOCK_USER_LINE_ID,
			DisplayName: cfg.MOCK_USER_NAME,
		}, nil
	}

	if err := h.verifyAccessToken(accessToken); err != nil {
		return nil, fmt.Errorf("failed to verify access token: %w", err)
	}
	req, err := http.NewRequest("GET", "https://api.line.me/v2/profile", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user: LINE API returned status %d: %s", resp.StatusCode, string(body))
	}

	var profile LineProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (h *registerHandler) verifyAccessToken(accessToken string) error {
	verifyURL := fmt.Sprintf("https://api.line.me/oauth2/v2.1/verify?access_token=%s", accessToken)

	resp, err := http.Get(verifyURL)
	if err != nil {
		return fmt.Errorf("failed to verify token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token verification failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		ClientID  string `json:"client_id"`
		ExpiresIn int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode verification response: %w", err)
	}

	// Verify that the token belongs to your LINE channel
	if result.ClientID != h.lineChannelID {
		return fmt.Errorf("user: token is not issued by the expected LINE channel (expected: %s, got: %s)",
			h.lineChannelID, result.ClientID)
	}

	return nil
}
