package auth

import "github.com/gin-gonic/gin"

type AuthHandler struct {
	svc AuthService
}

func NewAuthHandler(svc AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Signup(c *gin.Context)  {}
func (h *AuthHandler) Login(c *gin.Context)   {}
func (h *AuthHandler) Refresh(c *gin.Context) {}
func (h *AuthHandler) Logout(c *gin.Context)  {}
func (h *AuthHandler) GetMe(c *gin.Context)   {}
func (h *AuthHandler) UpdateMe(c *gin.Context) {}
