package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StaffOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || (role != "staff" && role != "admin") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé. Réservé au personnel et aux administrateurs."})
			c.Abort()
			return
		}
		c.Next()
	}
}