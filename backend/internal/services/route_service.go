package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"stops/backend/internal/models"
	"stops/backend/internal/repositories"
)

type RouteService struct {
	routes      *repositories.RouteRepository
	osrmBaseURL string
	httpClient  *http.Client
}

func NewRouteService(routes *repositories.RouteRepository, osrmBaseURL string) *RouteService {
	return &RouteService{
		routes:      routes,
		osrmBaseURL: osrmBaseURL,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

type OSRMResponse struct {
	Code    string      `json:"code"`
	Routes  []OSRMRoute `json:"routes"`
	Message string      `json:"message,omitempty"`
}

type OSRMRoute struct {
	Distance float64   `json:"distance"`
	Duration float64   `json:"duration"`
	Geometry string    `json:"geometry"`
	Legs     []OSRMLeg `json:"legs"`
}

type OSRMLeg struct {
	Distance float64    `json:"distance"`
	Duration float64    `json:"duration"`
	Steps    []OSRMStep `json:"steps"`
}

type OSRMStep struct {
	Instruction string  `json:"instruction"`
	Distance    float64 `json:"distance"`
	Duration    float64 `json:"duration"`
}

type DetailedRoute struct {
	ID               string           `json:"id"`
	RouteName        string           `json:"routeName"`
	StartPoint       string           `json:"startPoint"`
	DestinationPoint string           `json:"destinationPoint"`
	Fare             float64          `json:"fare"`
	City             string           `json:"city"`
	EstimatedTime    int              `json:"estimatedTime"`
	Status           string           `json:"status"`
	OSRMRoute        *OSRMRoute       `json:"osrmRoute,omitempty"`
	WalkingSegments  []WalkingSegment `json:"walkingSegments,omitempty"`
	Transfers        []Transfer       `json:"transfers,omitempty"`
}

type WalkingSegment struct {
	From     string     `json:"from"`
	To       string     `json:"to"`
	Distance float64    `json:"distance"`
	Duration float64    `json:"duration"`
	Steps    []OSRMStep `json:"steps"`
}

type Transfer struct {
	StopName string `json:"stopName"`
	Type     string `json:"type"`
}

func (service *RouteService) GetRoutes(ctx context.Context, start, dest string) ([]models.Route, error) {
	if start == "" && dest == "" {
		return service.routes.FindAll(ctx)
	}
	return service.routes.FindByStartAndDestination(ctx, start, dest)
}

func (service *RouteService) GetDetailedRoute(ctx context.Context, startLat, startLng, endLat, endLng float64, transportType string) (*DetailedRoute, error) {
	// Get OSRM route
	osrmURL := fmt.Sprintf("%s/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		service.osrmBaseURL, startLng, startLat, endLng, endLat)

	resp, err := service.httpClient.Get(osrmURL)
	if err != nil {
		return nil, fmt.Errorf("OSRM request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read OSRM response: %w", err)
	}

	var osrmResp OSRMResponse
	if err := json.Unmarshal(body, &osrmResp); err != nil {
		return nil, fmt.Errorf("parse OSRM response: %w", err)
	}

	if osrmResp.Code != "Ok" {
		return nil, fmt.Errorf("OSRM error: %s", osrmResp.Message)
	}

	if len(osrmResp.Routes) == 0 {
		return nil, fmt.Errorf("no routes found")
	}

	route := osrmResp.Routes[0]

	// Create detailed route
	detailedRoute := &DetailedRoute{
		RouteName:        fmt.Sprintf("Route from (%.4f, %.4f) to (%.4f, %.4f)", startLat, startLng, endLat, endLng),
		StartPoint:       fmt.Sprintf("%.4f, %.4f", startLat, startLng),
		DestinationPoint: fmt.Sprintf("%.4f, %.4f", endLat, endLng),
		City:             "Addis Ababa",
		EstimatedTime:    int(route.Duration / 60), // Convert to minutes
		Status:           "Active",
		OSRMRoute:        &route,
	}

	// Add walking segments (simplified - in production would calculate actual walking distances)
	if len(route.Legs) > 0 {
		for _, leg := range route.Legs {
			if len(leg.Steps) > 0 {
				detailedRoute.WalkingSegments = append(detailedRoute.WalkingSegments, WalkingSegment{
					From:     "Start",
					To:       "End",
					Distance: leg.Distance,
					Duration: leg.Duration,
					Steps:    leg.Steps,
				})
			}
		}
	}

	// Estimate fare based on distance
	if transportType == "bus" {
		detailedRoute.Fare = calculateBusFare(route.Distance)
	} else if transportType == "taxi" {
		detailedRoute.Fare = calculateTaxiFare(route.Distance)
	}

	return detailedRoute, nil
}

func calculateBusFare(distance float64) float64 {
	// Simple fare calculation: base fare + per km
	baseFare := 5.0 // ETB
	perKm := 1.5    // ETB per km
	return baseFare + (distance/1000)*perKm
}

func calculateTaxiFare(distance float64) float64 {
	// Simple fare calculation: base fare + per km
	baseFare := 15.0 // ETB
	perKm := 4.0     // ETB per km
	return baseFare + (distance/1000)*perKm
}

func (service *RouteService) GetAlternativeRoutes(ctx context.Context, startLat, startLng, endLat, endLng float64, congestionLevel string) ([]*DetailedRoute, error) {
	var routes []*DetailedRoute

	// Get primary route
	primaryRoute, err := service.GetDetailedRoute(ctx, startLat, startLng, endLat, endLng, "bus")
	if err != nil {
		return nil, err
	}
	routes = append(routes, primaryRoute)

	// If congestion is high, suggest alternative routes
	if congestionLevel == "High" {
		// Add taxi alternative
		taxiRoute, err := service.GetDetailedRoute(ctx, startLat, startLng, endLat, endLng, "taxi")
		if err == nil {
			taxiRoute.RouteName = "Taxi Alternative (Faster during congestion)"
			routes = append(routes, taxiRoute)
		}
	}

	return routes, nil
}

func (service *RouteService) GetFareEstimate(ctx context.Context, routeID string, transportType string) (float64, error) {
	fares, err := service.routes.GetFaresByRouteID(ctx, routeID)
	if err != nil {
		return 0, err
	}

	for _, f := range fares {
		if f.TransportType == transportType {
			return f.Amount, nil
		}
	}

	// Fallback/Default if no explicit transport type fare exists
	if len(fares) > 0 {
		return fares[0].Amount, nil
	}

	return 0, nil
}
