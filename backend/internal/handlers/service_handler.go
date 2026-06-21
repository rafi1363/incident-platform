package handlers

import (
	"net/http"

	"incident-platform/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	service *services.ServiceService
}

func NewServiceHandler(service *services.ServiceService) *ServiceHandler {
	return &ServiceHandler{service: service}
}

func (h *ServiceHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	data, err := h.service.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
