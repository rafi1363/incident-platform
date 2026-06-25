package main

import (
	"context"
	"log"
	"time"

	"incident-platform/backend/internal/database"
	"incident-platform/backend/internal/handlers"
	"incident-platform/backend/internal/repository"
	"incident-platform/backend/internal/services"

	"github.com/gin-contrib/cors"
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

	// ─── CORS middleware ─────────────────────────────────────────────
	// Tells the browser "it's OK for my frontend (different port) to call me."
	// MUST come before the routes so it wraps every request, including the
	// preflight OPTIONS requests the browser sends for PUT/DELETE/JSON.
	r.Use(cors.New(cors.Config{
		// AllowOrigins: which frontends may call us. In dev that's Vite's port.
		// In production you'd put your real domain here instead.
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
		// AllowMethods: which HTTP verbs the frontend may use.
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// AllowHeaders: which request headers the frontend may send.
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// ExposeHeaders: let the frontend read these response headers.
		ExposeHeaders: []string{"Content-Length"},
		// AllowCredentials: cookies/auth. False for now (no auth yet).
		AllowCredentials: false,
		// Preflight cache: browser can remember the CORS answer for 12h,
		// so it doesn't send an OPTIONS request before EVERY call.
		MaxAge: 12 * time.Hour,
	}))

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
