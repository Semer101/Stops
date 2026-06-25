package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"stops/backend/internal/services"
	"stops/backend/internal/utils"
)

type ChatController struct {
	chat *services.ChatService
}

type ChatRequest struct {
	Messages []services.ChatMessage `json:"messages" binding:"required"`
}

type ChatResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func NewChatController(chat *services.ChatService) *ChatController {
	return &ChatController{chat: chat}
}

func (controller *ChatController) SendMessage(context *gin.Context) {
	var request ChatRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		writeError(context, utils.ErrInvalidInput)
		return
	}

	reply, err := controller.chat.Chat(context.Request.Context(), request.Messages)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, ChatResponse{
		Success: true,
		Message: "Message sent successfully.",
		Data:    reply,
	})
}
