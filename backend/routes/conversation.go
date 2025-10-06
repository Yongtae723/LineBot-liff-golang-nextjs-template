package routes

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/linebot-liff-template/backend/config"
	"github.com/linebot-liff-template/backend/logic/conversation"
	"github.com/linebot-liff-template/go_pkg/llm"
	"github.com/linebot-liff-template/go_pkg/repository"
)

type ConversationHandler struct{}

func NewConversationHandler() *ConversationHandler {
	return &ConversationHandler{}
}

type ConversationRequest struct {
	Message string `json:"message" binding:"required"`
}

type ConversationResponse struct {
	Response string `json:"response"`
}

func (h *ConversationHandler) GetConversations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	cfg := config.Load()
	ctx := context.Background()
	geminiClient, err := llm.NewGoogleGemini(ctx, cfg.GEMINI_API_KEY, "gemini-1.5-flash")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	convRepo := repository.NewConversationRepo()
	userRepo := repository.NewUserRepo()
	handler := conversation.NewHandler(convRepo, userRepo, geminiClient)

	conversations, err := handler.GetHistory(c.Request.Context(), userID.(string), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, conversations)
}

func (h *ConversationHandler) PostConversation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req ConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := config.Load()
	ctx := context.Background()
	geminiClient, err := llm.NewGoogleGemini(ctx, cfg.GEMINI_API_KEY, "gemini-2.5-flash-lite")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	convRepo := repository.NewConversationRepo()
	userRepo := repository.NewUserRepo()
	handler := conversation.NewHandler(convRepo, userRepo, geminiClient)

	response, err := handler.ProcessMessage(c.Request.Context(), userID.(string), req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ConversationResponse{Response: response})
}
