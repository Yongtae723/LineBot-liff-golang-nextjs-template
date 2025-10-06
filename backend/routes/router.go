package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/linebot-liff-template/backend/middleware"
)

func SetupRoutes(r *gin.Engine, userHandler *UserHandler, conversationHandler *ConversationHandler) {
	r.Use(middleware.CORS())

	api := r.Group("/api/v1")
	{
		api.POST("/user/register", userHandler.RegisterUser)

		authenticated := api.Group("")
		authenticated.Use(middleware.Auth())
		{
			authenticated.GET("/conversations", conversationHandler.GetConversations)
			authenticated.POST("/conversations", conversationHandler.PostConversation)
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
