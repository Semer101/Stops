package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"stops/backend/internal/config"
	"stops/backend/internal/controllers"
	"stops/backend/internal/middleware"
	"stops/backend/internal/models"
	"stops/backend/internal/repositories"
	"stops/backend/internal/services"
)

func NewRouter(cfg config.Config, db *pgxpool.Pool) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Configure CORS middleware for local frontend development
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.CORSOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	healthController := controllers.NewHealthController()
	userRepository := repositories.NewUserRepository(db)
	passwordService := services.NewPasswordService()
	tokenService := services.NewTokenService(cfg.JWTSecret, cfg.JWTTTL)
	authService := services.NewAuthService(userRepository, passwordService, tokenService)
	authController := controllers.NewAuthController(authService)
	profileService := services.NewProfileService(userRepository)
	profileController := controllers.NewProfileController(profileService)
	transitStopRepository := repositories.NewTransitStopRepository(db)
	transitStopService := services.NewTransitStopService(transitStopRepository)
	transitStopController := controllers.NewTransitStopController(transitStopService)

	routeRepository := repositories.NewRouteRepository(db)
	osrmBaseURL := "https://router.project-osrm.org"
	routeService := services.NewRouteService(routeRepository, osrmBaseURL)
	routeController := controllers.NewRouteController(routeService)

	crowdsourceRepository := repositories.NewCrowdsourcingRepository(db)
	crowdsourceService := services.NewCrowdsourcingService(crowdsourceRepository)
	crowdsourceController := controllers.NewCrowdsourceController(crowdsourceService)

	predictionService := services.NewPredictionService()
	predictionController := controllers.NewPredictionController(predictionService)

	chatService := services.NewChatService(cfg.GeminiAPIKey)
	chatController := controllers.NewChatController(chatService)

	adminRepository := repositories.NewAdminRepository(db)
	adminService := services.NewAdminService(adminRepository)
	adminController := controllers.NewAdminController(adminService)

	userPreferencesRepository := repositories.NewUserPreferencesRepository(db)
	userPreferencesService := services.NewUserPreferencesService(userPreferencesRepository, userRepository)
	userPreferencesController := controllers.NewUserPreferencesController(userPreferencesService)

	api := router.Group("/api")
	api.GET("/health", healthController.Show)

	auth := api.Group("/auth")
	auth.POST("/register", authController.Register)
	auth.POST("/login", authController.Login)
	auth.POST("/forgot-password", authController.ForgotPassword)
	auth.POST("/reset-password", authController.ResetPassword)

	profile := api.Group("/profile")
	profile.Use(middleware.RequireAuth(tokenService))
	profile.GET("/me", profileController.Show)
	profile.PUT("/preferences", userPreferencesController.UpdatePreferences)
	profile.GET("/saved-places", userPreferencesController.GetSavedPlaces)
	profile.POST("/saved-places", userPreferencesController.CreateSavedPlace)
	profile.DELETE("/saved-places/:id", userPreferencesController.DeleteSavedPlace)
	profile.GET("/trips", userPreferencesController.GetTripHistory)

	stops := api.Group("/stops")
	stops.GET("/nearby", transitStopController.Nearby)

	// Protected crowdsourcing routes
	protectedStops := api.Group("/stops")
	protectedStops.Use(middleware.RequireAuth(tokenService))
	protectedStops.POST("/report", crowdsourceController.Report)
	protectedStops.POST("/vote", crowdsourceController.Vote)

	api.GET("/routes", routeController.List)
	api.GET("/routes/detailed", routeController.DetailedRoute)
	api.GET("/routes/alternatives", routeController.AlternativeRoutes)
	api.GET("/fare", routeController.FareEstimate)
	api.GET("/predictions", predictionController.Predict)

	chat := api.Group("/chat")
	chat.Use(middleware.RequireAuth(tokenService))
	chat.POST("", chatController.SendMessage)

	// Protected admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.RequireAuth(tokenService), middleware.RequireRole(models.RoleAdmin))
	admin.GET("/reports", adminController.ListPending)
	admin.POST("/reports/:id/verify", adminController.Verify)
	admin.POST("/reports/:id/reject", adminController.Reject)
	admin.GET("/audit-logs", adminController.ListAuditLogs)
	admin.GET("/analytics", adminController.GetAnalytics)

	// Admin CRUD for stops
	admin.POST("/stops", adminController.CreateStop)
	admin.PUT("/stops/:id", adminController.UpdateStop)
	admin.DELETE("/stops/:id", adminController.DeleteStop)

	// Admin CRUD for routes
	admin.POST("/routes", adminController.CreateRoute)
	admin.PUT("/routes/:id", adminController.UpdateRoute)
	admin.DELETE("/routes/:id", adminController.DeleteRoute)

	// Admin CRUD for fares
	admin.POST("/fares", adminController.CreateFare)
	admin.PUT("/fares/:id", adminController.UpdateFare)
	admin.DELETE("/fares/:id", adminController.DeleteFare)

	return router
}
