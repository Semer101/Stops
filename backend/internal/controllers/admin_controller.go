package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"stops/backend/internal/middleware"
	"stops/backend/internal/services"
	"stops/backend/internal/utils"
)

type AdminController struct {
	admin *services.AdminService
}

type AdminActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewAdminController(admin *services.AdminService) *AdminController {
	return &AdminController{admin: admin}
}

func (controller *AdminController) ListPending(context *gin.Context) {
	reports, err := controller.admin.GetPendingReports(context.Request.Context())
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pending reports loaded successfully.",
		"data":    reports,
	})
}

func (controller *AdminController) Verify(context *gin.Context) {
	value, ok := context.Get(middleware.ContextUserID)
	if !ok {
		writeError(context, utils.ErrUnauthorized)
		return
	}
	adminUserID, ok := value.(string)
	if !ok || adminUserID == "" {
		writeError(context, utils.ErrUnauthorized)
		return
	}

	reportID := context.Param("id")
	if reportID == "" {
		writeError(context, utils.ErrInvalidInput)
		return
	}

	err := controller.admin.VerifyReport(context.Request.Context(), reportID, adminUserID)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, AdminActionResponse{
		Success: true,
		Message: "Report verified successfully by administrator.",
	})
}

func (controller *AdminController) Reject(context *gin.Context) {
	value, ok := context.Get(middleware.ContextUserID)
	if !ok {
		writeError(context, utils.ErrUnauthorized)
		return
	}
	adminUserID, ok := value.(string)
	if !ok || adminUserID == "" {
		writeError(context, utils.ErrUnauthorized)
		return
	}

	reportID := context.Param("id")
	if reportID == "" {
		writeError(context, utils.ErrInvalidInput)
		return
	}

	err := controller.admin.RejectReport(context.Request.Context(), reportID, adminUserID)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, AdminActionResponse{
		Success: true,
		Message: "Report rejected and stop deleted successfully by administrator.",
	})
}

func (controller *AdminController) ListAuditLogs(context *gin.Context) {
	logs, err := controller.admin.GetAuditLogs(context.Request.Context())
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Admin audit logs loaded successfully.",
		"data":    logs,
	})
}

func (controller *AdminController) CreateStop(context *gin.Context) {
	var req struct {
		Name     string  `json:"name" binding:"required"`
		Type     string  `json:"type" binding:"required"`
		Lat      float64 `json:"lat" binding:"required"`
		Lng      float64 `json:"lng" binding:"required"`
		Area     string  `json:"area"`
		Verified bool    `json:"verified"`
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		writeError(context, err)
		return
	}

	stopID, err := controller.admin.CreateStop(context.Request.Context(), req.Name, req.Type, req.Lat, req.Lng, req.Area, req.Verified)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Stop created successfully.",
		"data":    gin.H{"id": stopID},
	})
}

func (controller *AdminController) UpdateStop(context *gin.Context) {
	stopID := context.Param("id")
	if stopID == "" {
		writeError(context, utils.ErrInvalidInput)
		return
	}

	var req struct {
		Name     *string  `json:"name"`
		Type     *string  `json:"type"`
		Lat      *float64 `json:"lat"`
		Lng      *float64 `json:"lng"`
		Area     *string  `json:"area"`
		Verified *bool    `json:"verified"`
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		writeError(context, err)
		return
	}

	err := controller.admin.UpdateStop(context.Request.Context(), stopID, req.Name, req.Type, req.Lat, req.Lng, req.Area, req.Verified)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, AdminActionResponse{
		Success: true,
		Message: "Stop updated successfully.",
	})
}

func (controller *AdminController) DeleteStop(context *gin.Context) {
	stopID := context.Param("id")
	if stopID == "" {
		writeError(context, utils.ErrInvalidInput)
		return
	}

	err := controller.admin.DeleteStop(context.Request.Context(), stopID)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, AdminActionResponse{
		Success: true,
		Message: "Stop deleted successfully.",
	})
}

func (controller *AdminController) CreateRoute(context *gin.Context) {
	var req struct {
		RouteName        string  `json:"routeName" binding:"required"`
		StartPoint       string  `json:"startPoint" binding:"required"`
		DestinationPoint string  `json:"destinationPoint" binding:"required"`
		Fare             float64 `json:"fare"`
		City             string  `json:"city"`
		EstimatedTime    int     `json:"estimatedTime"`
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		writeError(context, err)
		return
	}

	routeID, err := controller.admin.CreateRoute(context.Request.Context(), req.RouteName, req.StartPoint, req.DestinationPoint, req.Fare, req.City, req.EstimatedTime)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Route created successfully.",
		"data":    gin.H{"id": routeID},
	})
}

func (controller *AdminController) UpdateRoute(context *gin.Context) {
	routeID := context.Param("id")
	if routeID == "" {
		writeError(context, utils.ErrInvalidInput)
		return
	}

	var req struct {
		RouteName        *string  `json:"routeName"`
		StartPoint       *string  `json:"startPoint"`
		DestinationPoint *string  `json:"destinationPoint"`
		Fare             *float64 `json:"fare"`
		City             *string  `json:"city"`
		EstimatedTime    *int     `json:"estimatedTime"`
		Status           *string  `json:"status"`
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		writeError(context, err)
		return
	}

	err := controller.admin.UpdateRoute(context.Request.Context(), routeID, req.RouteName, req.StartPoint, req.DestinationPoint, req.Fare, req.City, req.EstimatedTime, req.Status)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, AdminActionResponse{
		Success: true,
		Message: "Route updated successfully.",
	})
}

func (controller *AdminController) DeleteRoute(context *gin.Context) {
	routeID := context.Param("id")
	if routeID == "" {
		writeError(context, utils.ErrInvalidInput)
		return
	}

	err := controller.admin.DeleteRoute(context.Request.Context(), routeID)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, AdminActionResponse{
		Success: true,
		Message: "Route deleted successfully.",
	})
}

func (controller *AdminController) CreateFare(context *gin.Context) {
	var req struct {
		RouteID       string  `json:"routeId" binding:"required"`
		TransportType string  `json:"transportType" binding:"required"`
		Amount        float64 `json:"amount" binding:"required"`
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		writeError(context, err)
		return
	}

	fareID, err := controller.admin.CreateFare(context.Request.Context(), req.RouteID, req.TransportType, req.Amount)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Fare created successfully.",
		"data":    gin.H{"id": fareID},
	})
}

func (controller *AdminController) UpdateFare(context *gin.Context) {
	fareID := context.Param("id")
	if fareID == "" {
		writeError(context, utils.ErrInvalidInput)
		return
	}

	var req struct {
		Amount *float64 `json:"amount"`
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		writeError(context, err)
		return
	}

	err := controller.admin.UpdateFare(context.Request.Context(), fareID, req.Amount)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, AdminActionResponse{
		Success: true,
		Message: "Fare updated successfully.",
	})
}

func (controller *AdminController) DeleteFare(context *gin.Context) {
	fareID := context.Param("id")
	if fareID == "" {
		writeError(context, utils.ErrInvalidInput)
		return
	}

	err := controller.admin.DeleteFare(context.Request.Context(), fareID)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, AdminActionResponse{
		Success: true,
		Message: "Fare deleted successfully.",
	})
}

func (controller *AdminController) GetAnalytics(context *gin.Context) {
	analytics, err := controller.admin.GetAnalytics(context.Request.Context())
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Analytics retrieved successfully.",
		"data":    analytics,
	})
}
