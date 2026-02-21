package fridge

import "github.com/gin-gonic/gin"

type FridgeHandler struct {
	svc FridgeService
}

func NewFridgeHandler(svc FridgeService) *FridgeHandler {
	return &FridgeHandler{svc: svc}
}

func (h *FridgeHandler) List(c *gin.Context)        {}
func (h *FridgeHandler) Create(c *gin.Context)      {}
func (h *FridgeHandler) GetSummary(c *gin.Context)  {}
func (h *FridgeHandler) GetByID(c *gin.Context)     {}
func (h *FridgeHandler) Update(c *gin.Context)      {}
func (h *FridgeHandler) Delete(c *gin.Context)      {}
func (h *FridgeHandler) BulkDelete(c *gin.Context)  {}
