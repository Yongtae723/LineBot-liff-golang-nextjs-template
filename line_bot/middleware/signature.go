package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ValidateSignature(channelSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		signature := c.GetHeader("X-Line-Signature")
		if signature == "" {
			log.Error().Msg("Missing X-Line-Signature header")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing signature"})
			c.Abort()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read request body")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
			c.Abort()
			return
		}

		if !validateSignature(channelSecret, signature, body) {
			log.Error().Msg("Invalid signature")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
			c.Abort()
			return
		}

		c.Set("body", body)
		c.Next()
	}
}

func validateSignature(channelSecret, signature string, body []byte) bool {
	mac := hmac.New(sha256.New, []byte(channelSecret))
	mac.Write(body)
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
