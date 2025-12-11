package middleware

import (
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func CasbinMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("userID")
		org, _ := c.Get("orgID")
		userIDStr := user.(string)
		orgIDStr := org.(string)
		path := c.Request.URL.Path
		method := c.Request.Method

		if user == "" || org == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User dan Organization ID tidak valid"})
			c.Abort()
			return
		}

		allowed, err := enforcer.Enforce(userIDStr, orgIDStr, path, method)

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
				"detail": "Role anda di " + orgIDStr + " tidak mengizinkan " + method,
			})
			c.Abort()
		}
	}
}
