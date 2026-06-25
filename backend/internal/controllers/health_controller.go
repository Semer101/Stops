package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

type HealthResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (controller *HealthController) Show(context *gin.Context) {
	context.JSON(http.StatusOK, HealthResponse{
		Success:   true,
		Message:   "STOPS API is healthy",
		Service:   "stops-backend",
		Timestamp: time.Now().UTC(),
	})
}
