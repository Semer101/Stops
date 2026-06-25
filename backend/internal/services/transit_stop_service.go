package services

import (
	"context"
	"math"

	"stops/backend/internal/models"
	"stops/backend/internal/repositories"
	"stops/backend/internal/utils"
)

const (
	defaultNearbyRadiusMeters = 3000
	maxNearbyRadiusMeters     = 10000
	defaultNearbyLimit        = 20
	maxNearbyLimit            = 50
	walkingMetersPerMinute    = 80
)

type TransitStopStore interface {
	FindNearby(ctx context.Context, params repositories.NearbyStopsParams) ([]models.TransitStop, error)
}

type TransitStopService struct {
	stops TransitStopStore
}

type NearbyStopsInput struct {
	Latitude     float64
	Longitude    float64
	RadiusMeters int
	Limit        int
}

func NewTransitStopService(stops TransitStopStore) *TransitStopService {
	return &TransitStopService{stops: stops}
}

func (service *TransitStopService) FindNearby(ctx context.Context, input NearbyStopsInput) ([]models.PublicTransitStop, error) {
	params, err := normalizeNearbyStopsInput(input)
	if err != nil {
		return nil, err
	}

	stops, err := service.stops.FindNearby(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]models.PublicTransitStop, 0, len(stops))
	for _, stop := range stops {
		stop.WalkingTimeMinutes = estimateWalkingMinutes(stop.DistanceMeters)
		result = append(result, stop.Public())
	}

	return result, nil
}

func normalizeNearbyStopsInput(input NearbyStopsInput) (repositories.NearbyStopsParams, error) {
	if input.Latitude < -90 || input.Latitude > 90 || input.Longitude < -180 || input.Longitude > 180 {
		return repositories.NearbyStopsParams{}, utils.ErrInvalidInput
	}

	radius := input.RadiusMeters
	if radius == 0 {
		radius = defaultNearbyRadiusMeters
	}
	if radius < 1 || radius > maxNearbyRadiusMeters {
		return repositories.NearbyStopsParams{}, utils.ErrInvalidInput
	}

	limit := input.Limit
	if limit == 0 {
		limit = defaultNearbyLimit
	}
	if limit < 1 || limit > maxNearbyLimit {
		return repositories.NearbyStopsParams{}, utils.ErrInvalidInput
	}

	return repositories.NearbyStopsParams{
		Latitude:     input.Latitude,
		Longitude:    input.Longitude,
		RadiusMeters: radius,
		Limit:        limit,
	}, nil
}

func estimateWalkingMinutes(distanceMeters float64) int {
	if distanceMeters <= 0 {
		return 0
	}

	return int(math.Ceil(distanceMeters / walkingMetersPerMinute))
}
