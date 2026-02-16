package main

import (
	"net/http"

	"gored/services"
	"gored/utils/ratelimit"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(ratelimit.Middleware())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})

	router.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.POST("/task", func(c *gin.Context) {
		err := services.Publish(c, []byte("do work"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "producer failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "queued"})
	})

	router.Run(":8080")
}
