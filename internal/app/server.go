package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"html/template"
	"io/fs"
	"log"
	"music_app/internal/databases"
	"music_app/internal/middleware"
	"music_app/internal/music/handler"
	"music_app/internal/music/model"
	"music_app/internal/music/repository"
	"music_app/internal/music/service"
	"music_app/internal/user"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (a *App) Run() error {
	bucket, err := databases.NewGridFSBucket(a.mongoClient.Database)
	if err != nil {
		return fmt.Errorf("failed to create GridFS bucket: %w", err)
	}

	db := a.mongoClient.Database
	a.trackRepo = repository.NewTrackRepository(db.Collection("tracks"), bucket)
	a.trackService = service.NewTrackService(a.trackRepo)
	a.trackHandler = handler.NewTrackHandler(a.trackService)

	router := gin.Default()
	router.Static("/static", "./web/static")
	funcMap := getFuncMap()
	htmlTemplate := template.Must(template.New("").Funcs(funcMap).ParseGlob("web/templates/**/*.html"))
	router.SetHTMLTemplate(htmlTemplate)
	router.Use(middleware.TemplateVars(a.sessionManager, a.userRepo))

	svc := user.NewService(a.userRepo, a.sessionManager)
	userHandler := user.NewHandler(svc, a.jwtManager)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"IsAuthenticated": c.GetBool("IsAuthenticated"),
			"Email":           c.GetString("Email"),
		})
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"IsAuthenticated": c.GetBool("IsAuthenticated"),
		})
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"IsAuthenticated": c.GetBool("IsAuthenticated"),
		})
	})

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	a.SeedExampleTracks()

	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware(a.sessionManager))

	auth.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"IsAuthenticated": true,
			"Email":           c.GetString("Email"),
		})
	})

	auth.POST("/logout", userHandler.Logout)
	auth.POST("/tracks", a.trackHandler.UploadTrackHandler)

	auth.GET("/tracks", a.trackHandler.ListTracksHandler)

	auth.GET("/recommendations", a.trackHandler.RecommendationsHandler)

	auth.GET("/tracks/:id/playview", func(c *gin.Context) {
		id := c.Param("id")
		track, err := a.trackHandler.GetTrackData(c.Request.Context(), id)
		if err != nil {
			c.String(http.StatusNotFound, "not found")
			return
		}
		c.HTML(http.StatusOK, "play.html", gin.H{
			"Track":           track,
			"IsAuthenticated": true,
			"Email":           c.GetString("Email"),
		})
	})

	auth.GET("/tracks/:id", a.trackHandler.GetTrackHandler)
	auth.GET("/tracks/:id/play", a.trackHandler.PlayTrackHandler)
	auth.DELETE("/tracks/:id", a.trackHandler.DeleteTrackHandler)

	admin := auth.Group("/admin")
	admin.Use(middleware.RequireAdmin())
	admin.GET("/users", userHandler.ListUsers)
	admin.POST("/users/:id/role", userHandler.ChangeUserRole)

	log.Printf("Starting server on %s", a.cfg.App.HTTPPort)
	return router.Run(a.cfg.App.HTTPPort)
}

func (a *App) SeedExampleTracks() {
	ctx := context.Background()

	count, err := a.mongoClient.Database.Collection("tracks").CountDocuments(ctx, bson.D{})
	if err != nil || count > 0 {
		return
	}

	assetsPath := "assets"

	err = filepath.WalkDir(assetsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to read directory: %w", err)
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".mp3") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		trackFileName := strings.TrimSuffix(d.Name(), ".mp3")
		trackInfo := strings.Split(trackFileName, "-")
		if len(trackInfo) != 3 {
			log.Printf("skipping malformed file name: %s", d.Name())
			return nil
		}
		trackName := strings.Replace(strings.TrimSpace(trackInfo[1]), "_", " ", -1)
		artistName := strings.Replace(strings.TrimSpace(trackInfo[0]), "_", " ", -1)
		genreName := strings.Replace(strings.TrimSpace(trackInfo[2]), "_", " ", -1)

		err = a.trackService.UploadTrack(ctx, &model.Track{
			Title:  trackName,
			Artist: artistName,
			Genre:  genreName,
		}, file)

		if err != nil {
			return fmt.Errorf("failed to upload track: %w", err)
		}

		return nil
	})

	if err != nil {
		log.Printf("failed to seed example tracks: %v", err)
	}
}

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
	}
}
