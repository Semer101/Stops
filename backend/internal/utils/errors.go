package utils

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyInUse  = errors.New("email already in use")
	ErrPhoneAlreadyInUse  = errors.New("phone already in use")
	ErrIdentifierRequired = errors.New("email or phone is required")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrUserNotFound       = errors.New("user not found")
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func ErrorStatus(err error) int {
	switch {
	case errors.Is(err, ErrInvalidInput),
		errors.Is(err, ErrEmailAlreadyInUse),
		errors.Is(err, ErrPhoneAlreadyInUse),
		errors.Is(err, ErrIdentifierRequired):
		return http.StatusBadRequest
	case errors.Is(err, ErrInvalidCredentials),
		errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrUserNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func PublicMessage(err error) string {
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		return "Invalid email, phone, or password."
	case errors.Is(err, ErrUnauthorized):
		return "Authentication is required."
	case errors.Is(err, ErrForbidden):
		return "You do not have permission to perform this action."
	case errors.Is(err, ErrEmailAlreadyInUse):
		return "Email is already in use."
	case errors.Is(err, ErrPhoneAlreadyInUse):
		return "Phone number is already in use."
	case errors.Is(err, ErrIdentifierRequired):
		return "Email or phone is required."
	case errors.Is(err, ErrUserNotFound):
		return "User was not found."
	case errors.Is(err, ErrInvalidInput):
		return "Request data is invalid."
	default:
		return "An unexpected error occurred."
	}
}
