package middleware

import (
	"github.com/gin-gonic/gin"
	"music_app/internal/user"
	"music_app/pkg/auth"
)

func TemplateVars(sm auth.SessionManager, ur user.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			c.Set("IsAuthenticated", false)
			c.Next()
			return
		}

		userID, err := sm.GetUserID(c, cookie)
		if err != nil {
			c.Set("IsAuthenticated", false)
			c.Next()
			return
		}

		u, err := ur.GetByID(c, userID)
		if err != nil {
			c.Set("IsAuthenticated", false)
			c.Next()
			return
		}

		c.Set("user", u)
		c.Set("IsAuthenticated", true)
		c.Set("Email", u.Email)
		c.Set("Role", u.Role)
		c.Next()
	}
}
