package main

import (
	"log"

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

	r := gin.Default()

	r.GET("/api/services", handler.GetAll)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	log.Println("API running on: 8080")

	r.Run(":8080")
}
