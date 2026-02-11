package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, data any, statusCode int) {
	c.JSON(statusCode, gin.H{"data": data})
}

func Error(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, gin.H{"error": err.Error()})
}
