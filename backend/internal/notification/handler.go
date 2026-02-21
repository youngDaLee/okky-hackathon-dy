package notification

import "github.com/gin-gonic/gin"

type NotificationHandler struct {
	svc NotificationService
}

func NewNotificationHandler(svc NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) List(c *gin.Context)    {}
func (h *NotificationHandler) Count(c *gin.Context)   {}
func (h *NotificationHandler) Read(c *gin.Context)    {}
func (h *NotificationHandler) ReadAll(c *gin.Context) {}
