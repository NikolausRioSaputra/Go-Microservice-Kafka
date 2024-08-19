package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("Request URI: %s\n", c.Request.RequestURI)
		c.Next()
	}
}
