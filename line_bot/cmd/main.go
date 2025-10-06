package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/linebot-liff-template/go_pkg/llm"
	"github.com/linebot-liff-template/go_pkg/repository"
	"github.com/linebot-liff-template/line_bot/config"
	"github.com/linebot-liff-template/line_bot/logic"
	"github.com/linebot-liff-template/line_bot/routes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Caller().Logger()

	if os.Getenv("ENV") == "local" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

		if err := godotenv.Load(); err != nil {
			log.Warn().Msg("No .env file found")
		}
		gin.SetMode(gin.DebugMode)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}

	cfg := config.Load()

	if err := repository.InitSupabase(cfg.SUPABASE_URL, cfg.SUPABASE_KEY); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Supabase")
	}

	bot, err := messaging_api.NewMessagingApiAPI(cfg.LINE_CHANNEL_TOKEN)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create LINE bot client")
	}

	userRepo := repository.NewUserRepo()
	convRepo := repository.NewConversationRepo()

	ctx := context.Background()
	geminiClient, err := llm.NewGoogleGemini(ctx, cfg.GEMINI_API_KEY, "gemini-2.5-flash-lite")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Gemini client")
	}

	messageHandler := logic.NewMessageHandler(convRepo, userRepo, geminiClient)

	followHandler, err := logic.NewFollowHandler(cfg.LINE_CHANNEL_TOKEN, cfg.LIFF_APP_URL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create follow handler")
	}

	webhookHandler := routes.NewWebhookHandler(bot, messageHandler, followHandler)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	routes.SetupRoutes(r, webhookHandler, cfg.LINE_CHANNEL_SECRET)

	addr := fmt.Sprintf(":%s", cfg.PORT)
	log.Info().Msgf("Starting LINE Bot server on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
