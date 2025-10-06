package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Sub string `json:"sub"`
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := validateSupabaseAuth(c); err != nil {
			log.Error().Err(err).
				Str("severity", "ERROR").
				Msg("Failed to validate Supabase authentication")
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}

func validateSupabaseAuth(c *gin.Context) error {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return fmt.Errorf("authorization header is required")
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		return fmt.Errorf("invalid token format")
	}

	token, err := jwt.ParseWithClaims(bearerToken[1], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		jwtSecret := os.Getenv("SUPABASE_JWT_SECRET")
		if jwtSecret == "" {
			return nil, fmt.Errorf("JWT secret is not set")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if claims.Sub == "" {
			return fmt.Errorf("user ID not found")
		}
		c.Set("user_id", claims.Sub)
		return nil
	}

	return fmt.Errorf("invalid token claims")
}
