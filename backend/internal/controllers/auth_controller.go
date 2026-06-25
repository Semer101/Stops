package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"stops/backend/internal/services"
	"stops/backend/internal/utils"
)

type AuthController struct {
	auth *services.AuthService
}

type RegisterRequest struct {
	Name     string  `json:"name" binding:"required"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Password string  `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Password string  `json:"password" binding:"required"`
}

type ForgotPasswordRequest struct {
	EmailOrPhone string `json:"email_or_phone" binding:"required"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type AuthResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    services.AuthResult `json:"data"`
}

func NewAuthController(auth *services.AuthService) *AuthController {
	return &AuthController{auth: auth}
}

func (controller *AuthController) Register(context *gin.Context) {
	var request RegisterRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		writeError(context, utils.ErrInvalidInput)
		return
	}

	result, err := controller.auth.Register(context.Request.Context(), services.RegisterInput{
		Name:     request.Name,
		Email:    request.Email,
		Phone:    request.Phone,
		Password: request.Password,
	})
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusCreated, AuthResponse{
		Success: true,
		Message: "Account registered successfully.",
		Data:    result,
	})
}

func (controller *AuthController) Login(context *gin.Context) {
	var request LoginRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		writeError(context, utils.ErrInvalidInput)
		return
	}

	result, err := controller.auth.Login(context.Request.Context(), services.LoginInput{
		Email:    request.Email,
		Phone:    request.Phone,
		Password: request.Password,
	})
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, AuthResponse{
		Success: true,
		Message: "Login successful.",
		Data:    result,
	})
}

func (controller *AuthController) ForgotPassword(context *gin.Context) {
	var request ForgotPasswordRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		writeError(context, utils.ErrInvalidInput)
		return
	}

	token, err := controller.auth.ForgotPassword(context.Request.Context(), request.EmailOrPhone)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password reset token generated successfully (Simulation: see console logs).",
		"token":   token,
	})
}

func (controller *AuthController) ResetPassword(context *gin.Context) {
	var request ResetPasswordRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		writeError(context, utils.ErrInvalidInput)
		return
	}

	err := controller.auth.ResetPassword(context.Request.Context(), request.Token, request.NewPassword)
	if err != nil {
		writeError(context, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password has been reset successfully.",
	})
}

func writeError(context *gin.Context, err error) {
	context.JSON(utils.ErrorStatus(err), utils.ErrorResponse{
		Success: false,
		Message: utils.PublicMessage(err),
		Error:   err.Error(),
	})
}
