package services

import (
	"context"

	"stops/backend/internal/models"
	"stops/backend/internal/repositories"
)

type UserPreferencesService struct {
	preferencesRepo *repositories.UserPreferencesRepository
	userRepo         *repositories.UserRepository
}

func NewUserPreferencesService(preferencesRepo *repositories.UserPreferencesRepository, userRepo *repositories.UserRepository) *UserPreferencesService {
	return &UserPreferencesService{
		preferencesRepo: preferencesRepo,
		userRepo:         userRepo,
	}
}

func (service *UserPreferencesService) UpdatePreferences(ctx context.Context, userID string, homeLocation, workLocation, preferredTransport *string) error {
	return service.preferencesRepo.UpdatePreferences(ctx, userID, homeLocation, workLocation, preferredTransport)
}

func (service *UserPreferencesService) GetSavedPlaces(ctx context.Context, userID string) ([]models.SavedPlace, error) {
	return service.preferencesRepo.GetSavedPlaces(ctx, userID)
}

func (service *UserPreferencesService) CreateSavedPlace(ctx context.Context, place models.SavedPlace) error {
	return service.preferencesRepo.CreateSavedPlace(ctx, place)
}

func (service *UserPreferencesService) DeleteSavedPlace(ctx context.Context, userID, placeID string) error {
	return service.preferencesRepo.DeleteSavedPlace(ctx, userID, placeID)
}

func (service *UserPreferencesService) GetTripHistory(ctx context.Context, userID string, limit int) ([]models.Trip, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return service.preferencesRepo.GetTripHistory(ctx, userID, limit)
}

func (service *UserPreferencesService) RecordTrip(ctx context.Context, trip models.Trip) error {
	return service.preferencesRepo.RecordTrip(ctx, trip)
}
