package handler

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"music_app/internal/music/model"
	"music_app/internal/music/service"
)

type TrackHandler struct {
	svc service.TrackService
}

func NewTrackHandler(svc service.TrackService) *TrackHandler {
	return &TrackHandler{svc: svc}
}

func (h *TrackHandler) UploadTrackHandler(c *gin.Context) {
	title := c.PostForm("title")
	artist := c.PostForm("artist")
	genre := c.PostForm("genre")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	if header.Size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty file"})
		return
	}

	track := &model.Track{
		Title:  title,
		Artist: artist,
		Genre:  genre,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	err = h.svc.UploadTrack(ctx, track, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload track: " + err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/tracks")
}

func (h *TrackHandler) ListTracksHandlerData(ctx context.Context, query string, page, pageSize int) ([]*model.Track, bool, error) {
	return h.svc.ListTracks(ctx, query, page, pageSize)
}

func (h *TrackHandler) ListTracksHandler(c *gin.Context) {
	query := c.Query("q")
	pageStr := c.Query("page")
	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	tracks, hasNext, err := h.svc.ListTracks(ctx, query, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tracks"})
		return
	}

	c.HTML(http.StatusOK, "tracks.html", gin.H{
		"Tracks":          tracks,
		"IsAuthenticated": c.GetBool("IsAuthenticated"),
		"Email":           c.GetString("Email"),
		"Role":            c.GetString("Role"),
		"Page":            page,
		"HasNext":         hasNext,
		"Query":           query,
	})
}

func (h *TrackHandler) GetTrackHandler(c *gin.Context) {
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	track, err := h.svc.GetTrack(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
		return
	}

	c.JSON(http.StatusOK, track)
}

func (h *TrackHandler) DeleteTrackHandler(c *gin.Context) {
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := h.svc.DeleteTrack(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete track"})
		return
	}

	c.Status(http.StatusOK)
}

func (h *TrackHandler) PlayTrackHandler(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetInt64("UserID")

	track, err := h.svc.GetTrack(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
		return
	}

	h.svc.TrackListening(c.Request.Context(), userID, id)

	downloadStream, err := h.svc.DownloadTrack(track.FileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open download stream"})
		return
	}
	defer downloadStream.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, downloadStream)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	reader := bytes.NewReader(buf.Bytes())
	http.ServeContent(c.Writer, c.Request, track.Title+".mp3", time.Now(), reader)
}

func (h *TrackHandler) GetTrackData(ctx context.Context, id string) (*model.Track, error) {
	return h.svc.GetTrack(ctx, id)
}

func (h *TrackHandler) RecommendationsHandler(c *gin.Context) {
	userID := c.GetInt64("UserID")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	tracks, err := h.svc.GetRecommendations(ctx, userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Ошибка при получении рекомендаций"})
		return
	}

	c.HTML(http.StatusOK, "recommendations.html", gin.H{
		"Tracks":          tracks,
		"Email":           c.GetString("Email"),
		"Role":            c.GetString("Role"),
		"IsAuthenticated": true,
	})
}
