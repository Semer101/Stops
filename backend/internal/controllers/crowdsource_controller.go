package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"stops/backend/internal/middleware"
	"stops/backend/internal/services"
	"stops/backend/internal/utils"
)

type CrowdsourceController struct {
	crowdsource *services.CrowdsourcingService
}

type ReportStopRequest struct {
	Name      string  `json:"name" binding:"required"`
	Type      string  `json:"type" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Area      string  `json:"area"`
}

type VoteRequest struct {
	ReportID  string `json:"report_id" binding:"required"`
	VoteValue int    `json:"vote_value" binding:"required"`
}

type CrowdsourceResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewCrowdsourceController(crowdsource *services.CrowdsourcingService) *CrowdsourceController {
	return &CrowdsourceController{crowdsource: crowdsource}
}

func (controller *CrowdsourceController) Report(context *gin.Context) {
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

	var request ReportStopRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		writeError(context, utils.ErrInvalidInput)
		return
	}

	reportID, err := controller.crowdsource.ReportStop(
		context.Request.Context(),
		request.Name,
		request.Type,
		request.Latitude,
		request.Longitude,
		request.Area,
		userID,
	)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusCreated, CrowdsourceResponse{
		Success: true,
		Message: "Stop reported successfully.",
		Data: gin.H{
			"report_id": reportID,
		},
	})
}

func (controller *CrowdsourceController) Vote(context *gin.Context) {
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

	var request VoteRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		writeError(context, utils.ErrInvalidInput)
		return
	}

	netVotes, err := controller.crowdsource.VoteStop(
		context.Request.Context(),
		request.ReportID,
		userID,
		request.VoteValue,
	)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, CrowdsourceResponse{
		Success: true,
		Message: "Vote submitted successfully.",
		Data: gin.H{
			"net_votes": netVotes,
		},
	})
}
