package handlers

import (
	"net/http"
	"tienda-backend/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type AssistantHandler struct {
	service ports.AssistantService
}

func NewAssistantHandler(service ports.AssistantService) *AssistantHandler {
	return &AssistantHandler{
		service: service,
	}
}

func (h *AssistantHandler) HandleMessage(c *gin.Context) {
	var input struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message is required"})
		return
	}

	response, err := h.service.GetResponse(input.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get response from assistant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}
