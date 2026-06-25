package services

import (
	"context"

	"stops/backend/internal/models"
)

type UserReader interface {
	FindByID(ctx context.Context, id string) (models.User, error)
}

type ProfileService struct {
	users UserReader
}

func NewProfileService(users UserReader) *ProfileService {
	return &ProfileService{users: users}
}

func (service *ProfileService) GetCurrentUser(ctx context.Context, userID string) (models.PublicUser, error) {
	user, err := service.users.FindByID(ctx, userID)
	if err != nil {
		return models.PublicUser{}, err
	}

	return user.Public(), nil
}
