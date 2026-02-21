package cookbook

import "github.com/gin-gonic/gin"

type CookbookHandler struct {
	svc CookbookService
}

func NewCookbookHandler(svc CookbookService) *CookbookHandler {
	return &CookbookHandler{svc: svc}
}

func (h *CookbookHandler) List(c *gin.Context)   {}
func (h *CookbookHandler) Save(c *gin.Context)   {}
func (h *CookbookHandler) Labels(c *gin.Context) {}
func (h *CookbookHandler) GetByID(c *gin.Context) {}
func (h *CookbookHandler) Update(c *gin.Context) {}
func (h *CookbookHandler) Delete(c *gin.Context) {}
