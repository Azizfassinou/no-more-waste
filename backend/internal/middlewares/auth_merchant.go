package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MerchantOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "merchant" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Accès refusé. Vous devez être un commerçant."})
			c.Abort()
			return
		}
		c.Next()
	}
}
