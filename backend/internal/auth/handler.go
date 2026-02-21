package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc AuthService
}

func NewAuthHandler(svc AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "message": err.Error()})
		return
	}
	resp, err := h.svc.Signup(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "message": err.Error()})
		return
	}
	resp, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "message": err.Error()})
		return
	}
	resp, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "message": err.Error()})
		return
	}
	if err := h.svc.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := c.GetString("userID")
	resp, err := h.svc.GetMe(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) UpdateMe(c *gin.Context) {
	var req UpdateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "message": err.Error()})
		return
	}
	userID := c.GetString("userID")
	resp, err := h.svc.UpdateMe(c.Request.Context(), userID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrDuplicateEmail):
		c.JSON(http.StatusConflict, gin.H{"error": "DUPLICATE_EMAIL", "message": "이미 사용 중인 이메일입니다"})
	case errors.Is(err, ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "INVALID_CREDENTIALS", "message": "이메일 또는 비밀번호가 올바르지 않습니다"})
	case errors.Is(err, ErrInvalidRefreshToken):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "INVALID_REFRESH_TOKEN", "message": "유효하지 않은 Refresh Token입니다. 다시 로그인해주세요."})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "NOT_FOUND", "message": "사용자를 찾을 수 없습니다"})
	case errors.Is(err, ErrValidation):
		c.JSON(http.StatusBadRequest, gin.H{"error": "VALIDATION_ERROR", "message": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INTERNAL_SERVER_ERROR", "message": "서버 오류가 발생했습니다"})
	}
}
