package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"io"

	"music_app/internal/music/model"
	"music_app/internal/music/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrackService interface {
	UploadTrack(ctx context.Context, track *model.Track, file io.Reader) error
	GetTrack(ctx context.Context, id string) (*model.Track, error)
	ListTracks(ctx context.Context) ([]*model.Track, error)
	DeleteTrack(ctx context.Context, id string) error
	DownloadTrack(fileID primitive.ObjectID) (*gridfs.DownloadStream, error)
}

type trackService struct {
	repo repository.TrackRepository
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

func (s *trackService) ListTracks(ctx context.Context) ([]*model.Track, error) {
	return s.repo.ListTracks(ctx)
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
