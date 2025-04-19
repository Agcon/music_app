package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

	trackCollection := a.mongoClient.Database.Collection("tracks")
	a.trackRepo = repository.NewTrackRepository(trackCollection, bucket)
	a.trackService = service.NewTrackService(a.trackRepo)
	a.trackHandler = handler.NewTrackHandler(a.trackService)

	router := gin.Default()
	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("web/templates/*.html")
	router.Use(middleware.TemplateVars(a.sessionManager, a.userRepo))

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

	svc := user.NewService(a.userRepo, a.sessionManager)
	userHandler := user.NewHandler(svc, a.jwtManager)

	a.SeedExampleTracks()

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

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

	auth.GET("/tracks", func(c *gin.Context) {
		tracks, err := a.trackHandler.ListTracksHandlerData(c.Request.Context())
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to get tracks")
			return
		}
		c.HTML(http.StatusOK, "tracks.html", gin.H{
			"Tracks":          tracks,
			"IsAuthenticated": true,
			"Email":           c.GetString("Email"),
		})
	})

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
	counter := 1

	filepath.WalkDir(assetsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".mp3") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		trackName := strings.TrimSuffix(d.Name(), ".mp3")
		artist := fmt.Sprintf("Artist %d", counter)

		a.trackService.UploadTrack(ctx, &model.Track{
			Title:  trackName,
			Artist: strings.TrimSpace(artist),
			Genre:  fmt.Sprintf("Genre %d", counter),
		}, file)

		counter++
		return nil
	})
}
