package services

import (
	"testing"
	"time"

	"stops/backend/internal/models"
)

func TestTokenServiceIssueAndParse(t *testing.T) {
	service := NewTokenService("test-secret", time.Hour)
	user := models.User{
		ID:   "user-1",
		Role: models.RoleAdmin,
	}

	token, err := service.Issue(user)
	if err != nil {
		t.Fatalf("issue token failed: %v", err)
	}

	claims, err := service.Parse(token)
	if err != nil {
		t.Fatalf("parse token failed: %v", err)
	}

	if claims.UserID != user.ID {
		t.Fatalf("expected user id %s, got %s", user.ID, claims.UserID)
	}
	if claims.Role != user.Role {
		t.Fatalf("expected role %s, got %s", user.Role, claims.Role)
	}
}

func TestTokenServiceRejectsInvalidToken(t *testing.T) {
	service := NewTokenService("test-secret", time.Hour)

	_, err := service.Parse("not-a-token")
	if err == nil {
		t.Fatal("expected invalid token error")
	}
}
