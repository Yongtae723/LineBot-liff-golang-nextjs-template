package config

import "os"

type Config struct {
	ENV                 string
	PORT                string
	LINE_CHANNEL_SECRET string
	LINE_CHANNEL_TOKEN  string
	SUPABASE_URL        string
	SUPABASE_KEY        string
	GEMINI_API_KEY      string
	BACKEND_URL         string
	LIFF_APP_URL        string
}

func Load() *Config {
	return &Config{
		ENV:                 getEnv("ENV", "local"),
		PORT:                getEnv("PORT", "8081"),
		LINE_CHANNEL_SECRET: getEnv("LINE_CHANNEL_SECRET", ""),
		LINE_CHANNEL_TOKEN:  getEnv("LINE_CHANNEL_TOKEN", ""),
		SUPABASE_URL:        getEnv("SUPABASE_URL", ""),
		SUPABASE_KEY:        getEnv("SUPABASE_KEY", ""),
		GEMINI_API_KEY:      getEnv("GEMINI_API_KEY", ""),
		BACKEND_URL:         getEnv("BACKEND_URL", "http://localhost:8080"),
		LIFF_APP_URL:        getEnv("LIFF_APP_URL", "https://liff.line.me/your-liff-id"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
