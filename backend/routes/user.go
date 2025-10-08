package routes

import (
	"net/http"

	"cookforyou.com/linebot-liff-template/backend/config"
	"cookforyou.com/linebot-liff-template/backend/logic/user"
	"cookforyou.com/linebot-liff-template/common/repository"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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

func (h *UserHandler) RegisterUserFromLiff(c *gin.Context) {
	cfg := config.Load()
	userRepo := repository.NewUserRepo()
	authRepo := repository.NewAuthRepo()
	registerHandler := user.NewRegisterHandler(userRepo, authRepo, cfg.LINE_CHANNEL_ID)

	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lineID, err := registerHandler.RegisterFromAccessToken(c.Request.Context(), req.AccessToken)
	if err != nil {
		log.Error().Err(err).Str("access_token", req.AccessToken).Msg("Failed to register user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserRegisterResponse{LineID: lineID})
}

type BotUserRegisterRequest struct {
	LineID      string `json:"line_id" binding:"required"`
	DisplayName string `json:"display_name"`
}

func (h *UserHandler) RegisterBotUser(c *gin.Context) {
	cfg := config.Load()
	userRepo := repository.NewUserRepo()
	authRepo := repository.NewAuthRepo()
	registerHandler := user.NewRegisterHandler(userRepo, authRepo, cfg.LINE_CHANNEL_ID)

	var req BotUserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lineID, err := registerHandler.Register(c.Request.Context(), req.LineID, req.DisplayName)
	if err != nil {
		log.Error().Err(err).Str("line_id", req.LineID).Msg("Failed to register user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserRegisterResponse{LineID: lineID})
}
