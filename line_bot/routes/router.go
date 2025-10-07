package routes

import (
	"cookforyou.com/linebot-liff-template/line_bot/config"
	"cookforyou.com/linebot-liff-template/line_bot/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	cfg := config.Load()

	r.POST("/webhook", middleware.ValidateSignature(cfg.LINE_CHANNEL_SECRET), HandleWebhook)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
