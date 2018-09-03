package api

import (
	"github.com/gin-gonic/gin"
)

// NotFound is the default not found handler
func NotFound(c *gin.Context) {
	c.JSON(404, gin.H{"error": "Not Found"})
}
