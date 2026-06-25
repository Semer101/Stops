package services

import (
	"context"
	"errors"
	"testing"

	"stops/backend/internal/models"
	"stops/backend/internal/repositories"
	"stops/backend/internal/utils"
)

type fakeTransitStopStore struct {
	params repositories.NearbyStopsParams
	stops  []models.TransitStop
	err    error
}

func (store *fakeTransitStopStore) FindNearby(_ context.Context, params repositories.NearbyStopsParams) ([]models.TransitStop, error) {
	store.params = params
	if store.err != nil {
		return nil, store.err
	}

	return store.stops, nil
}

func TestTransitStopServiceDefaultsAndWalkingTime(t *testing.T) {
	store := &fakeTransitStopStore{
		stops: []models.TransitStop{
			{
				ID:                 "stop-1",
				Name:               "Meskel Square",
				Type:               models.TransitStopTypeBusStop,
				Latitude:           9.0108,
				Longitude:          38.7612,
				Verified:           true,
				AvailabilityStatus: "Unknown",
				DistanceMeters:     161,
			},
		},
	}
	service := NewTransitStopService(store)

	stops, err := service.FindNearby(context.Background(), NearbyStopsInput{
		Latitude:  9.01,
		Longitude: 38.76,
	})
	if err != nil {
		t.Fatalf("find nearby failed: %v", err)
	}

	if store.params.RadiusMeters != defaultNearbyRadiusMeters {
		t.Fatalf("expected default radius, got %d", store.params.RadiusMeters)
	}
	if store.params.Limit != defaultNearbyLimit {
		t.Fatalf("expected default limit, got %d", store.params.Limit)
	}
	if stops[0].WalkingTimeMinutes != 3 {
		t.Fatalf("expected walking time 3, got %d", stops[0].WalkingTimeMinutes)
	}
}

func TestTransitStopServiceRejectsInvalidCoordinates(t *testing.T) {
	service := NewTransitStopService(&fakeTransitStopStore{})

	_, err := service.FindNearby(context.Background(), NearbyStopsInput{
		Latitude:  100,
		Longitude: 38.76,
	})

	if !errors.Is(err, utils.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestTransitStopServiceRejectsTooLargeRadius(t *testing.T) {
	service := NewTransitStopService(&fakeTransitStopStore{})

	_, err := service.FindNearby(context.Background(), NearbyStopsInput{
		Latitude:     9.01,
		Longitude:    38.76,
		RadiusMeters: maxNearbyRadiusMeters + 1,
	})

	if !errors.Is(err, utils.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}
