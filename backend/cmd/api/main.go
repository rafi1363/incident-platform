package main

import (
	"context"
	"log"
	"time"

	"incident-platform/backend/internal/database"
	"incident-platform/backend/internal/handlers"
	"incident-platform/backend/internal/repository"
	"incident-platform/backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	pool, err := database.NewPostgresPool()
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewServiceRepository(pool)
	service := services.NewServiceService(repo)
	handler := handlers.NewServiceHandler(service)
	incRepo := repository.NewIncidentRepository(pool)
	incService := services.NewIncidentService(incRepo)
	incHandler := handlers.NewIncidentHandler(incService)

	monitor := services.NewMonitorService(service, incService, 30*time.Second)
	monitor.Start(context.Background())

	r := gin.Default()

	r.GET("/api/services", handler.GetAll)
	r.POST("/api/services", handler.Create)
	r.GET("/api/services/:id", handler.GetByID)
	r.PUT("/api/services/:id", handler.Update)
	r.DELETE("/api/services/:id", handler.Delete)
	r.GET("/api/incidents", incHandler.GetAll)
	r.POST("/api/incidents", incHandler.Create)
	r.GET("/api/incidents/:id", incHandler.GetByID)
	r.GET("/api/services/:id/incidents", incHandler.GetByServiceID) // nested route
	r.PUT("/api/incidents/:id", incHandler.Update)
	r.DELETE("/api/incidents/:id", incHandler.Delete)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	log.Println("API running on: 8080")

	r.Run(":8080")
}
