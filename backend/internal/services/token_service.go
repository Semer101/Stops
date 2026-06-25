package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"stops/backend/internal/models"
	"stops/backend/internal/utils"
)

type TokenService struct {
	secret []byte
	ttl    time.Duration
}

type AuthClaims struct {
	UserID string          `json:"userId"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func NewTokenService(secret string, ttl time.Duration) *TokenService {
	return &TokenService{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (service *TokenService) Issue(user models.User) (string, error) {
	now := time.Now().UTC()
	claims := AuthClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(service.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(service.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signed, nil
}

func (service *TokenService) Parse(rawToken string) (AuthClaims, error) {
	token, err := jwt.ParseWithClaims(rawToken, &AuthClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, utils.ErrUnauthorized
		}

		return service.secret, nil
	})
	if err != nil {
		return AuthClaims{}, utils.ErrUnauthorized
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		return AuthClaims{}, utils.ErrUnauthorized
	}

	return *claims, nil
}
