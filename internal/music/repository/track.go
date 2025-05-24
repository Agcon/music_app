package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"time"

	"music_app/internal/music/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type TrackRepository interface {
	UploadTrack(ctx context.Context, track *model.Track, file io.Reader) error
	GetTrack(ctx context.Context, id primitive.ObjectID) (*model.Track, error)
	ListTracks(ctx context.Context) ([]*model.Track, error)
	DeleteTrack(ctx context.Context, id primitive.ObjectID) error
	ListTracksPaginated(ctx context.Context, query string, page, pageSize int) ([]*model.Track, bool, error)
	DownloadStreamFile(fileID primitive.ObjectID) (*gridfs.DownloadStream, error)
	FindByGenre(ctx context.Context, genre string, limit int) ([]*model.Track, error)
}

type trackRepository struct {
	collection *mongo.Collection
	bucket     *gridfs.Bucket
}

func NewTrackRepository(collection *mongo.Collection, bucket *gridfs.Bucket) TrackRepository {
	return &trackRepository{
		collection: collection,
		bucket:     bucket,
	}
}

func (r *trackRepository) UploadTrack(ctx context.Context, track *model.Track, file io.Reader) error {
	uploadStream, err := r.bucket.OpenUploadStream(
		track.Title,
	)
	if err != nil {
		return fmt.Errorf("failed to open upload stream: %w", err)
	}
	defer uploadStream.Close()

	_, err = io.Copy(uploadStream, file)
	if err != nil {
		return fmt.Errorf("failed to upload file to gridfs: %w", err)
	}

	track.FileID = uploadStream.FileID.(primitive.ObjectID)
	track.UploadedAt = time.Now()

	_, err = r.collection.InsertOne(ctx, track)
	if err != nil {
		return fmt.Errorf("failed to insert track metadata: %w", err)
	}

	return nil
}

func (r *trackRepository) GetTrack(ctx context.Context, id primitive.ObjectID) (*model.Track, error) {
	var track model.Track
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&track)
	if err != nil {
		return nil, fmt.Errorf("failed to find track: %w", err)
	}
	return &track, nil
}

func (r *trackRepository) ListTracks(ctx context.Context) ([]*model.Track, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to find tracks: %w", err)
	}
	defer cursor.Close(ctx)

	var tracks []*model.Track
	for cursor.Next(ctx) {
		var t model.Track
		if err = cursor.Decode(&t); err != nil {
			return nil, fmt.Errorf("failed to decode track: %w", err)
		}
		tracks = append(tracks, &t)
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate cursor: %w", err)
	}

	return tracks, nil
}

func (r *trackRepository) DeleteTrack(ctx context.Context, id primitive.ObjectID) error {
	track, err := r.GetTrack(ctx, id)
	if err != nil {
		return err
	}

	err = r.bucket.Delete(track.FileID)
	if err != nil {
		return fmt.Errorf("failed to delete file from gridfs: %w", err)
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete track metadata: %w", err)
	}

	return nil
}

func (r *trackRepository) DownloadStreamFile(fileID primitive.ObjectID) (*gridfs.DownloadStream, error) {
	return r.bucket.OpenDownloadStream(fileID)
}

func (r *trackRepository) ListTracksPaginated(ctx context.Context, query string, page, pageSize int) ([]*model.Track, bool, error) {
	filter := bson.M{}

	if query != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"title": bson.M{"$regex": query, "$options": "i"}},
				{"artist": bson.M{"$regex": query, "$options": "i"}},
			},
		}
	}

	skip := (page - 1) * pageSize

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize + 1)).
		SetSort(bson.D{{Key: "uploaded_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, false, fmt.Errorf("failed to find tracks: %w", err)
	}
	defer cursor.Close(ctx)

	var tracks []*model.Track
	for cursor.Next(ctx) {
		var t model.Track
		if err = cursor.Decode(&t); err != nil {
			return nil, false, fmt.Errorf("failed to decode track: %w", err)
		}
		tracks = append(tracks, &t)
	}

	hasNext := false
	if len(tracks) > pageSize {
		hasNext = true
		tracks = tracks[:pageSize]
	}

	if err = cursor.Err(); err != nil {
		return nil, false, fmt.Errorf("cursor error: %w", err)
	}

	return tracks, hasNext, nil
}

func (r *trackRepository) FindByGenre(ctx context.Context, genre string, limit int) ([]*model.Track, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"genre": genre}, options.Find().SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tracks []*model.Track
	for cursor.Next(ctx) {
		var t model.Track
		if err = cursor.Decode(&t); err != nil {
			continue
		}
		tracks = append(tracks, &t)
	}
	return tracks, nil
}
