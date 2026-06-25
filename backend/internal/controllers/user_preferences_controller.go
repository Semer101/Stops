package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"stops/backend/internal/models"
	"stops/backend/internal/services"
)

type UserPreferencesController struct {
	preferencesService *services.UserPreferencesService
}

func NewUserPreferencesController(preferencesService *services.UserPreferencesService) *UserPreferencesController {
	return &UserPreferencesController{preferencesService: preferencesService}
}

type PreferencesResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (controller *UserPreferencesController) UpdatePreferences(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		HomeLocation           *string `json:"homeLocation"`
		WorkLocation           *string `json:"workLocation"`
		PreferredTransportType *string `json:"preferredTransportType"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, PreferencesResponse{
			Success: false,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
		return
	}

	if err := controller.preferencesService.UpdatePreferences(c.Request.Context(), userID, req.HomeLocation, req.WorkLocation, req.PreferredTransportType); err != nil {
		c.JSON(http.StatusInternalServerError, PreferencesResponse{
			Success: false,
			Message: "Failed to update preferences",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PreferencesResponse{
		Success: true,
		Message: "Preferences updated successfully",
		Data:    nil,
	})
}

func (controller *UserPreferencesController) GetSavedPlaces(c *gin.Context) {
	userID := c.GetString("user_id")

	places, err := controller.preferencesService.GetSavedPlaces(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, PreferencesResponse{
			Success: false,
			Message: "Failed to get saved places",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PreferencesResponse{
		Success: true,
		Message: "Saved places retrieved successfully",
		Data:    places,
	})
}

func (controller *UserPreferencesController) CreateSavedPlace(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.SavedPlace
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, PreferencesResponse{
			Success: false,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
		return
	}

	req.UserID = userID

	if err := controller.preferencesService.CreateSavedPlace(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, PreferencesResponse{
			Success: false,
			Message: "Failed to create saved place",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, PreferencesResponse{
		Success: true,
		Message: "Saved place created successfully",
		Data:    req,
	})
}

func (controller *UserPreferencesController) DeleteSavedPlace(c *gin.Context) {
	userID := c.GetString("user_id")
	placeID := c.Param("id")

	if placeID == "" {
		c.JSON(http.StatusBadRequest, PreferencesResponse{
			Success: false,
			Message: "place_id is required",
		})
		return
	}

	if err := controller.preferencesService.DeleteSavedPlace(c.Request.Context(), userID, placeID); err != nil {
		c.JSON(http.StatusInternalServerError, PreferencesResponse{
			Success: false,
			Message: "Failed to delete saved place",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PreferencesResponse{
		Success: true,
		Message: "Saved place deleted successfully",
		Data:    nil,
	})
}

func (controller *UserPreferencesController) GetTripHistory(c *gin.Context) {
	userID := c.GetString("user_id")

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := parseInt(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	trips, err := controller.preferencesService.GetTripHistory(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, PreferencesResponse{
			Success: false,
			Message: "Failed to get trip history",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PreferencesResponse{
		Success: true,
		Message: "Trip history retrieved successfully",
		Data:    trips,
	})
}

func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}
