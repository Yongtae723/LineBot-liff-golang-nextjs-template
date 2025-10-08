package config

import "os"

type Config struct {
	ENV                 string
	PORT                string
	SUPABASE_URL        string
	SUPABASE_KEY        string
	SUPABASE_JWT_SECRET string
	GEMINI_API_KEY      string
	LINE_CHANNEL_ID     string
	MOCK_USER_LINE_ID   string
	MOCK_USER_NAME      string
}

func Load() *Config {
	return &Config{
		ENV:                 getEnv("ENV", "local"),
		PORT:                getEnv("PORT", "8080"),
		SUPABASE_URL:        getEnv("SUPABASE_URL", ""),
		SUPABASE_KEY:        getEnv("SUPABASE_KEY", ""),
		SUPABASE_JWT_SECRET: getEnv("SUPABASE_JWT_SECRET", ""),
		GEMINI_API_KEY:      getEnv("GEMINI_API_KEY", ""),
		LINE_CHANNEL_ID:     getEnv("LINE_CHANNEL_ID", ""),
		MOCK_USER_LINE_ID:   getEnv("MOCK_USER_LINE_ID", ""),
		MOCK_USER_NAME:      getEnv("MOCK_USER_NAME", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
