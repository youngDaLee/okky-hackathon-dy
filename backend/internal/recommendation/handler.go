package recommendation

import "github.com/gin-gonic/gin"

type RecommendationHandler struct {
	svc RecommendationService
}

func NewRecommendationHandler(svc RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{svc: svc}
}

func (h *RecommendationHandler) GetRecommendations(c *gin.Context)      {}
func (h *RecommendationHandler) GetTodayRecommendations(c *gin.Context) {}
func (h *RecommendationHandler) SearchRecipes(c *gin.Context)           {}
func (h *RecommendationHandler) GetRecipeByID(c *gin.Context)           {}
