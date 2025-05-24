package service

import (
	"context"
	"fmt"
	"io"
	"music_app/internal/music/model"
	"music_app/internal/music/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type TrackService interface {
	UploadTrack(ctx context.Context, track *model.Track, file io.Reader) error
	GetTrack(ctx context.Context, id string) (*model.Track, error)
	ListTracks(ctx context.Context, query string, page, pageSize int) ([]*model.Track, bool, error)
	DeleteTrack(ctx context.Context, id string) error
	DownloadTrack(fileID primitive.ObjectID) (*gridfs.DownloadStream, error)
	GetRecommendations(ctx context.Context, userID int64) ([]*model.Track, error)
	TrackListening(ctx context.Context, userID int64, trackIDHex string)
}

type trackService struct {
	repo        repository.TrackRepository
	historyRepo repository.HistoryRepository
}

func NewTrackService(repo repository.TrackRepository) TrackService {
	return &trackService{repo: repo}
}

func (s *trackService) UploadTrack(ctx context.Context, track *model.Track, file io.Reader) error {
	return s.repo.UploadTrack(ctx, track, file)
}

func (s *trackService) GetTrack(ctx context.Context, id string) (*model.Track, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse track id: %w", err)
	}
	return s.repo.GetTrack(ctx, objID)
}

func (s *trackService) ListTracks(ctx context.Context, query string, page, pageSize int) ([]*model.Track, bool, error) {
	return s.repo.ListTracksPaginated(ctx, query, page, pageSize)
}

func (s *trackService) DeleteTrack(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to parse track id: %w", err)
	}
	return s.repo.DeleteTrack(ctx, objID)
}

func (s *trackService) DownloadTrack(fileID primitive.ObjectID) (*gridfs.DownloadStream, error) {
	return s.repo.DownloadStreamFile(fileID)
}

func (s *trackService) GetRecommendations(ctx context.Context, userID int64) ([]*model.Track, error) {
	topGenres, err := s.historyRepo.GetTopGenres(ctx, userID)
	if err != nil || len(topGenres) == 0 {
		return nil, nil
	}
	return s.repo.FindByGenre(ctx, topGenres[0], 10)
}

func (s *trackService) TrackListening(ctx context.Context, userID int64, trackIDHex string) {
	trackID, err := primitive.ObjectIDFromHex(trackIDHex)
	if err != nil {
		return
	}

	track, err := s.repo.GetTrack(ctx, trackID)
	if err != nil {
		return
	}

	listen := model.TrackListen{
		UserID:    userID,
		TrackID:   track.ID,
		Genre:     track.Genre,
		Timestamp: time.Now(),
	}

	_ = s.repo.SaveTrackListen(ctx, listen)
}
