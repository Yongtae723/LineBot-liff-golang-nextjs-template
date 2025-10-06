package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/linebot-liff-template/line_bot/middleware"
)

func SetupRoutes(r *gin.Engine, webhookHandler *WebhookHandler, channelSecret string) {
	r.POST("/webhook", middleware.ValidateSignature(channelSecret), webhookHandler.Handle)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
