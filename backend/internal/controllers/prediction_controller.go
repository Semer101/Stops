package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"stops/backend/internal/services"
)

type PredictionController struct {
	predictions *services.PredictionService
}

type PredictionsResponse struct {
	Success bool                        `json:"success"`
	Message string                      `json:"message"`
	Data    []services.PredictionResult `json:"data"`
}

func NewPredictionController(predictions *services.PredictionService) *PredictionController {
	return &PredictionController{predictions: predictions}
}

func (controller *PredictionController) Predict(context *gin.Context) {
	routeID := context.Query("route_id")

	results, err := controller.predictions.GetPredictions(context.Request.Context(), routeID)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, PredictionsResponse{
		Success: true,
		Message: "Predictions generated successfully.",
		Data:    results,
	})
}
