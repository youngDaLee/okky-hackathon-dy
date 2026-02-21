package vision

import "github.com/gin-gonic/gin"

type VisionHandler struct {
	svc VisionService
}

func NewVisionHandler(svc VisionService) *VisionHandler {
	return &VisionHandler{svc: svc}
}

func (h *VisionHandler) CreateJob(c *gin.Context)  {}
func (h *VisionHandler) GetJob(c *gin.Context)     {}
func (h *VisionHandler) ConfirmJob(c *gin.Context) {}
