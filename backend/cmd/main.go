package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/linebot-liff-template/backend/config"
	"github.com/linebot-liff-template/backend/routes"
	"github.com/linebot-liff-template/go_pkg/repository"
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

	userHandler := routes.NewUserHandler()
	conversationHandler := routes.NewConversationHandler()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	routes.SetupRoutes(r, userHandler, conversationHandler)

	addr := fmt.Sprintf(":%s", cfg.PORT)
	log.Info().Msgf("Starting server on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
