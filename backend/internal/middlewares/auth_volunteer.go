package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VolunteerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || (role != "volunteer" && role != "staff" && role != "admin") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Accès refusé. Réservé aux bénévoles, staff ou administrateurs."})
			c.Abort()
			return
		}
		c.Next()
	}
}