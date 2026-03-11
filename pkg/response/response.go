package response

import "github.com/gin-gonic/gin"

func JSON(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"message": message,
		"data":    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"message": message,
	})
}