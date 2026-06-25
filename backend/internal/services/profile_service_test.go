package services

import (
	"context"
	"testing"
	"time"

	"stops/backend/internal/models"
)

type fakeUserReader struct {
	user models.User
	err  error
}

func (reader fakeUserReader) FindByID(_ context.Context, _ string) (models.User, error) {
	if reader.err != nil {
		return models.User{}, reader.err
	}

	return reader.user, nil
}

func TestProfileServiceReturnsPublicUser(t *testing.T) {
	email := "user@example.com"
	service := NewProfileService(fakeUserReader{
		user: models.User{
			ID:                "user-1",
			Name:              "Test User",
			Email:             &email,
			PasswordHash:      "secret",
			Role:              models.RoleUser,
			ContributionScore: 3,
			CreatedAt:         time.Now().UTC(),
		},
	})

	user, err := service.GetCurrentUser(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("get current user failed: %v", err)
	}

	if user.ID != "user-1" {
		t.Fatalf("expected user id, got %s", user.ID)
	}
	if user.Email == nil || *user.Email != email {
		t.Fatalf("expected email %s, got %#v", email, user.Email)
	}
}
