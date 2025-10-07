package auth

import (
	"time"

	"cookforyou.com/linebot-liff-template/line_bot/config"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWTToken(userID string) string {
	cfg := config.Load()

	claims := jwt.MapClaims{
		"aud":  "authenticated",
		"sub":  userID,
		"role": "authenticated",
		"app_metadata": map[string]string{
			"provider": "line",
			"userId":   userID,
		},
		"exp": time.Now().Add(time.Hour).Unix(),
		"iss": "supabase",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(cfg.SUPABASE_JWT_SECRET))
	return signedToken
}
