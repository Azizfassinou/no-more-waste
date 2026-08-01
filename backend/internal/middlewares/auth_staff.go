package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StaffOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "staff" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Accès refusé. Vous devez être un membre du personnel."})
			c.Abort()
			return
		}
		c.Next()
	}
}
