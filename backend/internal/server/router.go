package server

import (
	"github.com/gin-gonic/gin"
	"okky-hackathon/fridge-master-backend/internal/fridge"
	"okky-hackathon/fridge-master-backend/internal/server/middleware"
)

type RouterDeps struct {
	FridgeHandler *fridge.FridgeHandler
}

func NewRouter(deps RouterDeps) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api/v1")
	api.Use(middleware.Auth())

	// Task 5.1 — fridge 라우트 7개 등록
	f := deps.FridgeHandler
	fridgeGroup := api.Group("/fridge")
	{
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
