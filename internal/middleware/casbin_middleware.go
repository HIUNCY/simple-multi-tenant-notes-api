package middleware

import (
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func CasbinMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetHeader("X-User-ID")
		org := c.GetHeader("X-Organization-ID")
		path := c.Request.URL.Path
		method := c.Request.Method

		if user == "" || org == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User/Org Header missing"})
			c.Abort()
			return
		}

		allowed, err := enforcer.Enforce(user, org, path, method)

		if err != nil {
			log.Printf("Casbin error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authorization error"})
			c.Abort()
			return
		}

		if allowed {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "Anda tidak memiliki akses (Forbidden)",
				"detail": "Role anda di " + org + " tidak mengizinkan " + method,
			})
			c.Abort()
		}
	}
}
