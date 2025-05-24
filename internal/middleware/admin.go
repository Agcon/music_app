package middleware

import (
	"github.com/gin-gonic/gin"
	user2 "music_app/internal/user"
	"net/http"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, exists := c.Get("user")
		if !exists {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		userPtr, ok := u.(*user2.User)
		if !ok || userPtr.Role != "admin" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
