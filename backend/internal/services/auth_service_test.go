package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"stops/backend/internal/models"
	"stops/backend/internal/repositories"
	"stops/backend/internal/utils"
)

type fakeUserStore struct {
	user      models.User
	createErr error
	findErr   error
	created   repositories.CreateUserParams
}

func (store *fakeUserStore) Create(_ context.Context, params repositories.CreateUserParams) (models.User, error) {
	store.created = params
	if store.createErr != nil {
		return models.User{}, store.createErr
	}

	user := store.user
	user.Name = params.Name
	user.Email = params.Email
	user.Phone = params.Phone
	user.PasswordHash = params.PasswordHash
	user.Role = params.Role
	return user, nil
}

func (store *fakeUserStore) FindByEmailOrPhone(_ context.Context, _ *string, _ *string) (models.User, error) {
	if store.findErr != nil {
		return models.User{}, store.findErr
	}

	return store.user, nil
}

func (store *fakeUserStore) UpdateResetToken(_ context.Context, _ string, _ string, _ time.Time) error {
	return nil
}

func (store *fakeUserStore) FindByResetToken(_ context.Context, _ string) (models.User, error) {
	if store.findErr != nil {
		return models.User{}, store.findErr
	}
	return store.user, nil
}

func (store *fakeUserStore) UpdatePassword(_ context.Context, _ string, _ string) error {
	return nil
}

func TestRegisterRequiresEmailOrPhone(t *testing.T) {
	service := newTestAuthService(&fakeUserStore{})

	_, err := service.Register(context.Background(), RegisterInput{
		Name:     "Test User",
		Password: "password123",
	})

	if !errors.Is(err, utils.ErrIdentifierRequired) {
		t.Fatalf("expected identifier error, got %v", err)
	}
}

func TestRegisterHashesPasswordAndIssuesToken(t *testing.T) {
	email := "USER@Example.COM"
	store := &fakeUserStore{
		user: models.User{
			ID:        "user-1",
			CreatedAt: time.Now().UTC(),
		},
	}
	service := newTestAuthService(store)

	result, err := service.Register(context.Background(), RegisterInput{
		Name:     " Test User ",
		Email:    &email,
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	if result.Token == "" {
		t.Fatal("expected token")
	}
	if store.created.PasswordHash == "password123" {
		t.Fatal("expected password to be hashed")
	}
	if store.created.Email == nil || *store.created.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %#v", store.created.Email)
	}
	if result.User.Role != models.RoleUser {
		t.Fatalf("expected user role, got %s", result.User.Role)
	}
}

func TestLoginMapsMissingUserToInvalidCredentials(t *testing.T) {
	email := "user@example.com"
	store := &fakeUserStore{findErr: utils.ErrUserNotFound}
	service := newTestAuthService(store)

	_, err := service.Login(context.Background(), LoginInput{
		Email:    &email,
		Password: "password123",
	})

	if !errors.Is(err, utils.ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestLoginIssuesTokenForValidPassword(t *testing.T) {
	email := "user@example.com"
	passwords := NewPasswordService()
	hash, err := passwords.Hash("password123")
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}

	store := &fakeUserStore{
		user: models.User{
			ID:           "user-1",
			Name:         "Test User",
			Email:        &email,
			PasswordHash: hash,
			Role:         models.RoleUser,
			CreatedAt:    time.Now().UTC(),
		},
	}
	service := NewAuthService(store, passwords, NewTokenService("test-secret", time.Hour))

	result, err := service.Login(context.Background(), LoginInput{
		Email:    &email,
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	if result.Token == "" {
		t.Fatal("expected token")
	}
	if result.User.ID != "user-1" {
		t.Fatalf("expected public user, got %s", result.User.ID)
	}
}

func newTestAuthService(store *fakeUserStore) *AuthService {
	return NewAuthService(store, NewPasswordService(), NewTokenService("test-secret", time.Hour))
}
