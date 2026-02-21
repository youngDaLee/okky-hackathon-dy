package server

import (
	"github.com/gin-gonic/gin"
	"okky-hackathon/fridge-master-backend/internal/auth"
	"okky-hackathon/fridge-master-backend/internal/fridge"
	"okky-hackathon/fridge-master-backend/internal/server/middleware"
)

type RouterDeps struct {
	AuthHandler   *auth.AuthHandler
	AuthService   auth.AuthService
	FridgeHandler *fridge.FridgeHandler
}

func NewRouter(deps RouterDeps) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api/v1")

	// Public auth routes (no JWT required)
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/signup", deps.AuthHandler.Signup)
		authGroup.POST("/login", deps.AuthHandler.Login)
		authGroup.POST("/refresh", deps.AuthHandler.Refresh)
		authGroup.POST("/logout", deps.AuthHandler.Logout)
	}

	// Protected routes (JWT required)
	protected := api.Group("")
	protected.Use(middleware.Auth(deps.AuthService))
	{
		// Users
		users := protected.Group("/users")
		users.GET("/me", deps.AuthHandler.GetMe)
		users.PATCH("/me", deps.AuthHandler.UpdateMe)

		// Fridge
		f := deps.FridgeHandler
		fridgeGroup := protected.Group("/fridge")
		fridgeGroup.GET("", f.List)
		fridgeGroup.POST("", f.Create)
		fridgeGroup.DELETE("", f.BulkDelete)
		fridgeGroup.GET("/summary", f.GetSummary)
		fridgeGroup.GET("/:id", f.GetByID)
		fridgeGroup.PATCH("/:id", f.Update)
		fridgeGroup.DELETE("/:id", f.Delete)
	}

	return r
}
