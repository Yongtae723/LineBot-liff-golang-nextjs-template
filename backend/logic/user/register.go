package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/linebot-liff-template/go_pkg/models"
	"github.com/linebot-liff-template/go_pkg/repository"
)

type RegisterHandler struct {
	userRepo      repository.UserRepo
	authRepo      repository.AuthRepo
	lineChannelID string
}

func NewRegisterHandler(userRepo repository.UserRepo, authRepo repository.AuthRepo, lineChannelID string) *RegisterHandler {
	return &RegisterHandler{
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

func (h *RegisterHandler) Register(ctx context.Context, accessToken string) (string, error) {
	lineProfile, err := h.getLineProfile(accessToken)
	if err != nil {
		return "", fmt.Errorf("failed to get LINE profile: %w", err)
	}

	user, err := h.userRepo.GetByLineID(ctx, lineProfile.UserID)
	if err == nil {
		return user.LineID, nil
	}

	// NOTE: This is a simplified authentication flow for template purposes.
	// In production, implement proper authentication with:
	// - LINE Login integration with Supabase Auth
	// - Proper OAuth flow
	// - Secure token exchange
	// - Email verification (if needed)
	email := lineProfile.UserID
	password := lineProfile.UserID

	userMetadata := map[string]interface{}{
		"line_id": lineProfile.UserID,
		"name":    lineProfile.DisplayName,
	}
	appMetadata := map[string]interface{}{
		"provider": "line",
	}

	authUserID, err := h.authRepo.CreateUser(ctx, email, password, userMetadata, appMetadata)
	if err != nil && err != repository.ErrUserAlreadyExists {
		return "", fmt.Errorf("failed to create auth user: %w", err)
	}

	if err == repository.ErrUserAlreadyExists {
		authUserID, err = h.authRepo.GetUserIDByLineID(ctx, lineProfile.UserID)
		if err != nil {
			return "", fmt.Errorf("failed to get existing user: %w", err)
		}
	}

	newUser := &models.User{
		ID:     authUserID,
		LineID: lineProfile.UserID,
		Name:   lineProfile.DisplayName,
	}

	if err := h.userRepo.Create(ctx, newUser); err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return lineProfile.UserID, nil
}

func (h *RegisterHandler) getLineProfile(accessToken string) (*LineProfile, error) {
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
		return nil, fmt.Errorf("LINE API returned status %d: %s", resp.StatusCode, string(body))
	}

	var profile LineProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (h *RegisterHandler) verifyAccessToken(accessToken string) error {
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
		return fmt.Errorf("token is not issued by the expected LINE channel (expected: %s, got: %s)",
			h.lineChannelID, result.ClientID)
	}

	return nil
}
