package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"incident-platform/backend/internal/models"
	"incident-platform/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type IncidentHandler struct {
	incident *services.IncidentService
}

func NewIncidentHandler(incident *services.IncidentService) *IncidentHandler {
	return &IncidentHandler{incident: incident}
}

func (h *IncidentHandler) GetAll(c *gin.Context) {
	data, err := h.incident.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// Create reads JSON body into CreateIncidentInput. ShouldBindJSON does two jobs:
//  1. parse the JSON,  2. run binding:"required" validation.
func (h *IncidentHandler) Create(c *gin.Context) {
	var in models.CreateIncidentInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inc, err := h.incident.Create(c.Request.Context(), in)
	if err != nil {
		if err.Error() == "invalid incident status" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 201 Created is the correct status for a newly created resource.
	c.JSON(http.StatusCreated, inc)
}

// GetByID reads :id from the path. c.Param always returns a string, so we
// convert with strconv and reject garbage (e.g. /api/incidents/abc).
func (h *IncidentHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	inc, err := h.incident.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "incident not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inc)
}

// GetByServiceID — all incidents for a given service.
// Here :id from the path is the SERVICE id.
func (h *IncidentHandler) GetByServiceID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	data, err := h.incident.GetByServiceID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *IncidentHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var in models.UpdateIncidentInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inc, err := h.incident.Update(c.Request.Context(), id, in)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "incident not found"})
			return
		}
		if err.Error() == "invalid incident status" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inc)
}

func (h *IncidentHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.incident.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 204 No Content = success, and nothing to show for it.
	c.Status(http.StatusNoContent)
}
