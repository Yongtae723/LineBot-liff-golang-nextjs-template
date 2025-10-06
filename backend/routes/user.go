package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linebot-liff-template/backend/config"
	"github.com/linebot-liff-template/backend/logic/user"
	"github.com/linebot-liff-template/go_pkg/repository"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

type UserRegisterRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

type UserRegisterResponse struct {
	LineID string `json:"line_id"`
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	cfg := config.Load()
	userRepo := repository.NewUserRepo()
	authRepo := repository.NewAuthRepo()
	registerHandler := user.NewRegisterHandler(userRepo, authRepo, cfg.LINE_CHANNEL_ID)

	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lineID, err := registerHandler.Register(c.Request.Context(), req.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserRegisterResponse{LineID: lineID})
}
