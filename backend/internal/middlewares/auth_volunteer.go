package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VolunteerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "volunteer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Accès refusé. Vous devez être un bénévole."})
			c.Abort()
			return
		}
		c.Next()
	}
}
