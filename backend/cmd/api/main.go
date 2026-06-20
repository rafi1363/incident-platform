package main

import (
	"log"
	"net/http"

	"incident-platform/backend/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	_, err := database.NewPostgresPool()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	log.Println("API running on :8080 - main.go:26")

	r.Run(":8080")
}
