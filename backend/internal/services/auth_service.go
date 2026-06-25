package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"stops/backend/internal/models"
	"stops/backend/internal/repositories"
	"stops/backend/internal/utils"
)

type UserStore interface {
	Create(ctx context.Context, params repositories.CreateUserParams) (models.User, error)
	FindByEmailOrPhone(ctx context.Context, email *string, phone *string) (models.User, error)
	UpdateResetToken(ctx context.Context, userID string, token string, expiry time.Time) error
	FindByResetToken(ctx context.Context, token string) (models.User, error)
	UpdatePassword(ctx context.Context, userID string, passwordHash string) error
}

type AuthService struct {
	users     UserStore
	passwords *PasswordService
	tokens    *TokenService
}

type RegisterInput struct {
	Name     string
	Email    *string
	Phone    *string
	Password string
}

type LoginInput struct {
	Email    *string
	Phone    *string
	Password string
}

type AuthResult struct {
	User  models.PublicUser `json:"user"`
	Token string            `json:"token"`
}

func NewAuthService(users UserStore, passwords *PasswordService, tokens *TokenService) *AuthService {
	return &AuthService{
		users:     users,
		passwords: passwords,
		tokens:    tokens,
	}
}

func (service *AuthService) Register(ctx context.Context, input RegisterInput) (AuthResult, error) {
	normalized := normalizeRegisterInput(input)
	if err := validateRegisterInput(normalized); err != nil {
		return AuthResult{}, err
	}

	hash, err := service.passwords.Hash(normalized.Password)
	if err != nil {
		return AuthResult{}, err
	}

	user, err := service.users.Create(ctx, repositories.CreateUserParams{
		Name:         normalized.Name,
		Email:        normalized.Email,
		Phone:        normalized.Phone,
		PasswordHash: hash,
		Role:         models.RoleUser,
	})
	if err != nil {
		return AuthResult{}, err
	}

	token, err := service.tokens.Issue(user)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{User: user.Public(), Token: token}, nil
}

func (service *AuthService) Login(ctx context.Context, input LoginInput) (AuthResult, error) {
	normalized := normalizeLoginInput(input)
	if err := validateLoginInput(normalized); err != nil {
		return AuthResult{}, err
	}

	user, err := service.users.FindByEmailOrPhone(ctx, normalized.Email, normalized.Phone)
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			return AuthResult{}, utils.ErrInvalidCredentials
		}

		return AuthResult{}, err
	}

	if !service.passwords.Compare(user.PasswordHash, normalized.Password) {
		return AuthResult{}, utils.ErrInvalidCredentials
	}

	token, err := service.tokens.Issue(user)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{User: user.Public(), Token: token}, nil
}

func normalizeRegisterInput(input RegisterInput) RegisterInput {
	input.Name = strings.TrimSpace(input.Name)
	input.Password = strings.TrimSpace(input.Password)
	input.Email = normalizeOptional(input.Email, true)
	input.Phone = normalizeOptional(input.Phone, false)
	return input
}

func normalizeLoginInput(input LoginInput) LoginInput {
	input.Password = strings.TrimSpace(input.Password)
	input.Email = normalizeOptional(input.Email, true)
	input.Phone = normalizeOptional(input.Phone, false)
	return input
}

func normalizeOptional(value *string, lower bool) *string {
	if value == nil {
		return nil
	}

	normalized := strings.TrimSpace(*value)
	if lower {
		normalized = strings.ToLower(normalized)
	}

	if normalized == "" {
		return nil
	}

	return &normalized
}

func validateRegisterInput(input RegisterInput) error {
	if input.Name == "" || len(input.Password) < 8 {
		return utils.ErrInvalidInput
	}

	if input.Email == nil && input.Phone == nil {
		return utils.ErrIdentifierRequired
	}

	return nil
}

func validateLoginInput(input LoginInput) error {
	if input.Password == "" {
		return utils.ErrInvalidInput
	}

	if input.Email == nil && input.Phone == nil {
		return utils.ErrIdentifierRequired
	}

	return nil
}

func (service *AuthService) ForgotPassword(ctx context.Context, emailOrPhone string) (string, error) {
	emailOrPhone = strings.TrimSpace(emailOrPhone)
	if emailOrPhone == "" {
		return "", utils.ErrInvalidInput
	}

	var email, phone *string
	if strings.Contains(emailOrPhone, "@") {
		lower := strings.ToLower(emailOrPhone)
		email = &lower
	} else {
		phone = &emailOrPhone
	}

	user, err := service.users.FindByEmailOrPhone(ctx, email, phone)
	if err != nil {
		return "", err
	}

	// Generate secure 16-byte hex token
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	token := hex.EncodeToString(bytes)

	// Token expires in 1 hour
	expiry := time.Now().Add(1 * time.Hour)

	if err := service.users.UpdateResetToken(ctx, user.ID, token, expiry); err != nil {
		return "", err
	}

	// Local simulation logging
	log.Printf("[SIMULATION] Password reset token generated for user %s: %s", user.Name, token)

	return token, nil
}

func (service *AuthService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	token = strings.TrimSpace(token)
	newPassword = strings.TrimSpace(newPassword)
	if token == "" || len(newPassword) < 8 {
		return utils.ErrInvalidInput
	}

	user, err := service.users.FindByResetToken(ctx, token)
	if err != nil {
		return err
	}

	hash, err := service.passwords.Hash(newPassword)
	if err != nil {
		return err
	}

	return service.users.UpdatePassword(ctx, user.ID, hash)
}
