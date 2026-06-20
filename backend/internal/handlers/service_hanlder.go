package handlers

import (
	"net/http"

	"incident-platform/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	Repo *repository.ServiceRepository
}

func NewServiceHandler(
	repo *repository.ServiceRepository,
) *ServiceHandler {
	return &ServiceHandler{
		Repo: repo,
	}
}

func (h *ServiceHandler) GetServices(
	c *gin.Context,
) {
	services, err := h.Repo.GetAll(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		services,
	)
}
