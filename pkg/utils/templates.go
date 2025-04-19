package utils

import "github.com/gin-gonic/gin"

func Render(c *gin.Context, code int, name string, data gin.H) {
	if data == nil {
		data = gin.H{}
	}
	if isAuth, exists := c.Get("IsAuthenticated"); exists {
		data["IsAuthenticated"] = isAuth
	}
	if email, exists := c.Get("Email"); exists {
		data["Email"] = email
	}
	c.HTML(code, name, data)
}
