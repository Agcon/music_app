package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"music_app/pkg/auth"
	"net/http"
)

type Handler struct {
	svc        Service
	jwtManager auth.JWTManager
}

func NewHandler(svc Service, jwtManager auth.JWTManager) *Handler {
	return &Handler{svc: svc, jwtManager: jwtManager}
}

func (h *Handler) Register(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	if username == "" || email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Все поля обязательны"})
		return
	}

	err := h.svc.Register(context.Background(), username, email, password)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}

func (h *Handler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	if email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Все поля обязательны"})
		return
	}

	userID, sessionToken, err := h.svc.Login(context.Background(), email, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": err.Error()})
		return
	}

	jwtToken, err := h.jwtManager.Generate(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "не удалось создать токен"})
		return
	}

	log.Printf("ttl: %v", int(h.jwtManager.GetTTL()))

	c.SetCookie("session_token", sessionToken, int(h.jwtManager.GetTTL()), "/", "", false, true)
	c.SetCookie("jwt", jwtToken, int(h.jwtManager.GetTTL()), "/", "", false, true)

	c.Redirect(http.StatusSeeOther, "/tracks")
}

func (h *Handler) Logout(c *gin.Context) {
	token, err := c.Cookie("session_token")
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "отсутствует cookie сессии"})
		return
	}
	if err = h.svc.Logout(context.Background(), token); err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": "не удалось разлогиниться"})
		return
	}
	c.SetCookie("session_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}
