package utility

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    data,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"message": message,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": message,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{
		"message": message,
	})
}

func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": message,
	})
}
