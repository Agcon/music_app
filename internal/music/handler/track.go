package handler

import (
	"bytes"
	"context"
	"io"
	"net/http"
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

func (h *TrackHandler) ListTracksHandlerData(ctx context.Context) ([]*model.Track, error) {
	return h.svc.ListTracks(ctx)
}

func (h *TrackHandler) ListTracksHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	tracks, err := h.svc.ListTracks(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tracks"})
		return
	}

	c.JSON(http.StatusOK, tracks)
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
	track, err := h.svc.GetTrack(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "track not found"})
		return
	}

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
