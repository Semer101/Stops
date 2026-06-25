package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"stops/backend/internal/middleware"
	"stops/backend/internal/models"
	"stops/backend/internal/services"
	"stops/backend/internal/utils"
)

type ProfileController struct {
	profiles *services.ProfileService
}

type ProfileResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    models.PublicUser `json:"data"`
}

func NewProfileController(profiles *services.ProfileService) *ProfileController {
	return &ProfileController{profiles: profiles}
}

func (controller *ProfileController) Show(context *gin.Context) {
	value, ok := context.Get(middleware.ContextUserID)
	if !ok {
		writeError(context, utils.ErrUnauthorized)
		return
	}

	userID, ok := value.(string)
	if !ok || userID == "" {
		writeError(context, utils.ErrUnauthorized)
		return
	}

	user, err := controller.profiles.GetCurrentUser(context.Request.Context(), userID)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, ProfileResponse{
		Success: true,
		Message: "Profile loaded successfully.",
		Data:    user,
	})
}
