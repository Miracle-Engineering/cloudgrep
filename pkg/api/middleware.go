package api

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Middleware to print out request parameters and body for debugging
func requestInspectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.Request.ParseForm()
		log.Println("Request params:", err, c.Request.Form)
	}
}
