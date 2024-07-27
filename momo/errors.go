package momo

import (
	"log"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, statusCode int, err interface{}) {
	log.Printf("Error: %v", err)
	c.JSON(statusCode, gin.H{"error": err})
}
