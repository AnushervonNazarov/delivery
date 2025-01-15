package controllers

import (
	service "delivery/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDirections(c *gin.Context) {
	origin := c.Query("origin")
	destination := c.Query("destination")

	if origin == "" || destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "origin and destination are required"})
		return
	}

	// Call the function to get directions
	directions, err := service.GetDirections(origin, destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"directions": directions})
}
