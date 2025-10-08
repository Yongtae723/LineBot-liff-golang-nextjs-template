package routes

import (
	"cookforyou.com/linebot-liff-template/backend/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userHandler *UserHandler, conversationHandler *ConversationHandler) {
	r.Use(middleware.CORS())

	api := r.Group("/api/v1")

	api.POST("/user/register/liff", userHandler.RegisterUserFromLiff)

	authenticated := api.Group("")
	authenticated.Use(middleware.Auth())
	authenticated.POST("/user/register/bot", userHandler.RegisterBotUser)
	authenticated.GET("/conversations", conversationHandler.GetConversations)
	authenticated.POST("/conversations", conversationHandler.PostConversation)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
