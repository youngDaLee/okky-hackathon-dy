package recommendation

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
	svc RecommendationService
}

func NewRecommendationHandler(svc RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{svc: svc}
}

func getUserID(c *gin.Context) (string, bool) {
	v, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED"})
		return "", false
	}
	id, ok := v.(string)
	if !ok || id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED"})
		return "", false
	}
	return id, true
}

// Task 3.1 — GET /recommendations
func (h *RecommendationHandler) GetRecommendations(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	req := RecommendationRequest{
		Category: c.Query("category"),
		Limit:    20,
	}

	if v := c.Query("tier"); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil && n >= 1 && n <= 3 {
			req.Tier = &n
		}
	}
	if v := c.Query("max_missing"); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			req.MaxMissing = &n
		}
	}
	if v := c.Query("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil && n > 0 {
			req.Limit = n
		}
	}

	results, err := h.svc.GetRecommendations(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR"})
		return
	}

	if results == nil {
		results = []RecommendationResult{}
	}
	c.JSON(http.StatusOK, gin.H{
		"items": results,
		"total": len(results),
	})
}

// Task 3.2 — GET /recommendations/today
func (h *RecommendationHandler) GetTodayRecommendations(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	results, err := h.svc.GetTodayRecommendations(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR"})
		return
	}

	if results == nil {
		results = []RecommendationResult{}
	}
	c.JSON(http.StatusOK, gin.H{
		"items": results,
		"total": len(results),
	})
}

// Task 3.3 — GET /recipes
func (h *RecommendationHandler) SearchRecipes(c *gin.Context) {
	limit := 20
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}

	recipes, err := h.svc.SearchRecipes(
		c.Request.Context(),
		c.Query("keyword"),
		c.Query("category"),
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR"})
		return
	}

	if recipes == nil {
		recipes = []*Recipe{}
	}
	c.JSON(http.StatusOK, gin.H{
		"items": recipes,
		"total": len(recipes),
	})
}

// Task 3.4 — GET /recipes/:id
func (h *RecommendationHandler) GetRecipeByID(c *gin.Context) {
	recipe, err := h.svc.GetRecipeByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "RECIPE_NOT_FOUND"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR"})
		return
	}
	c.JSON(http.StatusOK, recipe)
}
