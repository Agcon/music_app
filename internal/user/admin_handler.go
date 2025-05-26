package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) AdminDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
		"IsAuthenticated": c.GetBool("IsAuthenticated"),
		"Email":           c.GetString("Email"),
		"Role":            c.GetString("Role"),
	})
}

func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.svc.ListAll(context.Background())
	if err != nil {
		c.String(http.StatusInternalServerError, "ошибка получения пользователей")
		return
	}

	c.HTML(http.StatusOK, "admin_users.html", gin.H{
		"Users":           users,
		"IsAuthenticated": c.GetBool("IsAuthenticated"),
		"Email":           c.GetString("Email"),
		"Role":            c.GetString("Role"),
	})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	intID, _ := strconv.ParseInt(id, 10, 64)
	if err := h.svc.DeleteUserByID(context.Background(), intID); err != nil {
		c.String(http.StatusInternalServerError, "не удалось удалить пользователя")
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/users")
}

func (h *Handler) ChangeUserRole(c *gin.Context) {
	userID := c.Param("id")
	intID, _ := strconv.ParseInt(userID, 10, 64)
	newRole := c.PostForm("role")
	err := h.svc.UpdateUserRole(c.Request.Context(), intID, newRole)
	if err != nil {
		c.String(http.StatusInternalServerError, "не удалось обновить роль")
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/users")
}
