package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"stops/backend/internal/models"
	"stops/backend/internal/services"
	"stops/backend/internal/utils"
)

type TransitStopController struct {
	stops *services.TransitStopService
}

type NearbyStopsResponse struct {
	Success bool                       `json:"success"`
	Message string                     `json:"message"`
	Data    []models.PublicTransitStop `json:"data"`
}

func NewTransitStopController(stops *services.TransitStopService) *TransitStopController {
	return &TransitStopController{stops: stops}
}

func (controller *TransitStopController) Nearby(context *gin.Context) {
	latitude, err := parseRequiredFloat(context, "lat")
	if err != nil {
		writeError(context, err)
		return
	}

	longitude, err := parseRequiredFloat(context, "lng")
	if err != nil {
		writeError(context, err)
		return
	}

	radius, err := parseOptionalInt(context, "radius")
	if err != nil {
		writeError(context, err)
		return
	}

	limit, err := parseOptionalInt(context, "limit")
	if err != nil {
		writeError(context, err)
		return
	}

	stops, err := controller.stops.FindNearby(context.Request.Context(), services.NearbyStopsInput{
		Latitude:     latitude,
		Longitude:    longitude,
		RadiusMeters: radius,
		Limit:        limit,
	})
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, NearbyStopsResponse{
		Success: true,
		Message: nearbyStopsMessage(stops),
		Data:    stops,
	})
}

func parseRequiredFloat(context *gin.Context, key string) (float64, error) {
	raw := context.Query(key)
	if raw == "" {
		return 0, utils.ErrInvalidInput
	}

	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, utils.ErrInvalidInput
	}

	return value, nil
}

func parseOptionalInt(context *gin.Context, key string) (int, error) {
	raw := context.Query(key)
	if raw == "" {
		return 0, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, utils.ErrInvalidInput
	}

	return value, nil
}

func nearbyStopsMessage(stops []models.PublicTransitStop) string {
	if len(stops) == 0 {
		return "This area hasn't been mapped yet. Help your community by adding the first transit stop."
	}

	return "Nearby transit stops loaded successfully."
}
