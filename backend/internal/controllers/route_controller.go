package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"stops/backend/internal/services"
)

type RouteController struct {
	routes *services.RouteService
}

func NewRouteController(routes *services.RouteService) *RouteController {
	return &RouteController{routes: routes}
}

type RoutesResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (controller *RouteController) List(context *gin.Context) {
	start := context.Query("start")
	dest := context.Query("dest")

	routes, err := controller.routes.GetRoutes(context.Request.Context(), start, dest)
	if err != nil {
		context.JSON(http.StatusInternalServerError, RoutesResponse{
			Success: false,
			Message: "Failed to query routes",
			Data:    err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, RoutesResponse{
		Success: true,
		Message: "Routes loaded successfully",
		Data:    routes,
	})
}

func (controller *RouteController) FareEstimate(context *gin.Context) {
	routeID := context.Query("route_id")
	transportType := context.Query("transport_type")

	if routeID == "" {
		context.JSON(http.StatusBadRequest, RoutesResponse{
			Success: false,
			Message: "route_id is required",
		})
		return
	}

	fare, err := controller.routes.GetFareEstimate(context.Request.Context(), routeID, transportType)
	if err != nil {
		context.JSON(http.StatusInternalServerError, RoutesResponse{
			Success: false,
			Message: "Failed to estimate fare",
			Data:    err.Error(),
		})
		return
	}

	if fare == 0 {
		context.JSON(http.StatusOK, RoutesResponse{
			Success: false,
			Message: "Fare data is missing for this route and transport type.",
			Data:    nil,
		})
		return
	}

	context.JSON(http.StatusOK, RoutesResponse{
		Success: true,
		Message: "Fare estimated successfully",
		Data: gin.H{
			"fare": fare,
		},
	})
}

func (controller *RouteController) DetailedRoute(context *gin.Context) {
	startLatStr := context.Query("start_lat")
	startLngStr := context.Query("start_lng")
	endLatStr := context.Query("end_lat")
	endLngStr := context.Query("end_lng")
	transportType := context.Query("transport_type")

	if startLatStr == "" || startLngStr == "" || endLatStr == "" || endLngStr == "" {
		context.JSON(http.StatusBadRequest, RoutesResponse{
			Success: false,
			Message: "start_lat, start_lng, end_lat, and end_lng are required",
		})
		return
	}

	startLat, _ := parseFloat(startLatStr)
	startLng, _ := parseFloat(startLngStr)
	endLat, _ := parseFloat(endLatStr)
	endLng, _ := parseFloat(endLngStr)

	if transportType == "" {
		transportType = "bus"
	}

	route, err := controller.routes.GetDetailedRoute(context.Request.Context(), startLat, startLng, endLat, endLng, transportType)
	if err != nil {
		context.JSON(http.StatusInternalServerError, RoutesResponse{
			Success: false,
			Message: "Failed to get detailed route",
			Data:    err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, RoutesResponse{
		Success: true,
		Message: "Detailed route retrieved successfully",
		Data:    route,
	})
}

func (controller *RouteController) AlternativeRoutes(context *gin.Context) {
	startLatStr := context.Query("start_lat")
	startLngStr := context.Query("start_lng")
	endLatStr := context.Query("end_lat")
	endLngStr := context.Query("end_lng")
	congestionLevel := context.Query("congestion_level")

	if startLatStr == "" || startLngStr == "" || endLatStr == "" || endLngStr == "" {
		context.JSON(http.StatusBadRequest, RoutesResponse{
			Success: false,
			Message: "start_lat, start_lng, end_lat, and end_lng are required",
		})
		return
	}

	startLat, _ := parseFloat(startLatStr)
	startLng, _ := parseFloat(startLngStr)
	endLat, _ := parseFloat(endLatStr)
	endLng, _ := parseFloat(endLngStr)

	if congestionLevel == "" {
		congestionLevel = "Medium"
	}

	routes, err := controller.routes.GetAlternativeRoutes(context.Request.Context(), startLat, startLng, endLat, endLng, congestionLevel)
	if err != nil {
		context.JSON(http.StatusInternalServerError, RoutesResponse{
			Success: false,
			Message: "Failed to get alternative routes",
			Data:    err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, RoutesResponse{
		Success: true,
		Message: "Alternative routes retrieved successfully",
		Data:    routes,
	})
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}
