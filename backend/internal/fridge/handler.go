package fridge

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FridgeHandler struct {
	svc FridgeService
}

func NewFridgeHandler(svc FridgeService) *FridgeHandler {
	return &FridgeHandler{svc: svc}
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

func handleServiceErr(c *gin.Context, err error) {
	var dupErr *DuplicateError
	switch {
	case errors.As(err, &dupErr):
		c.JSON(http.StatusConflict, gin.H{
			"error":       "DUPLICATE_INGREDIENT",
			"message":     "이미 등록된 재료입니다. 기존 수량에 합산하려면 PATCH를 사용하세요.",
			"existing_id": dupErr.ExistingID,
		})
	case errors.Is(err, ErrLimitExceeded):
		c.JSON(http.StatusBadRequest, gin.H{"error": "FRIDGE_LIMIT_EXCEEDED"})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "INGREDIENT_NOT_FOUND"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_ERROR"})
	}
}

// Task 4.1 — GET /fridge
func (h *FridgeHandler) List(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	filter := ListFilter{
		Category:     c.Query("category"),
		ExpiryStatus: c.Query("expiry_status"),
		Search:       c.Query("q"),
	}

	items, total, err := h.svc.ListIngredients(c.Request.Context(), userID, filter)
	if err != nil {
		handleServiceErr(c, err)
		return
	}

	if items == nil {
		items = []*IngredientResponse{}
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"total": total,
	})
}

// Task 4.2 — POST /fridge
func (h *FridgeHandler) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req CreateIngredientReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	if !isValidCategory(req.Category) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "message": "허용되지 않는 category"})
		return
	}

	resp, err := h.svc.AddIngredient(c.Request.Context(), userID, req)
	if err != nil {
		handleServiceErr(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Task 4.3 — GET /fridge/summary
func (h *FridgeHandler) GetSummary(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	summary, err := h.svc.GetSummary(c.Request.Context(), userID)
	if err != nil {
		handleServiceErr(c, err)
		return
	}

	c.JSON(http.StatusOK, summary)
}

// Task 4.4 — GET /fridge/:id
func (h *FridgeHandler) GetByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	resp, err := h.svc.GetIngredient(c.Request.Context(), c.Param("id"), userID)
	if err != nil {
		handleServiceErr(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Task 4.5 — PATCH /fridge/:id
func (h *FridgeHandler) Update(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req UpdateIngredientReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	if req.Category != nil && !isValidCategory(*req.Category) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "message": "허용되지 않는 category"})
		return
	}
	if req.Quantity != nil && *req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "message": "수량은 0보다 커야 합니다"})
		return
	}

	resp, err := h.svc.UpdateIngredient(c.Request.Context(), c.Param("id"), userID, req)
	if err != nil {
		handleServiceErr(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Task 4.6 — DELETE /fridge/:id
func (h *FridgeHandler) Delete(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	if err := h.svc.RemoveIngredient(c.Request.Context(), c.Param("id"), userID); err != nil {
		handleServiceErr(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Task 4.7 — DELETE /fridge (bulk)
func (h *FridgeHandler) BulkDelete(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req struct {
		IDs []string `json:"ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "message": err.Error()})
		return
	}

	count, err := h.svc.BulkRemove(c.Request.Context(), req.IDs, userID)
	if err != nil {
		handleServiceErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted_count": count})
}

func isValidCategory(cat string) bool {
	switch cat {
	case CategoryVegetable, CategoryFruit, CategoryMeat, CategorySeafood,
		CategoryDairy, CategoryGrain, CategoryCondiment, CategoryFrozen, CategoryOther:
		return true
	}
	return false
}
