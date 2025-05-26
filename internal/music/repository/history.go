package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"music_app/internal/music/model"
)

type HistoryRepository interface {
	SaveTrackListen(ctx context.Context, listen model.TrackListen) error
	GetTopGenres(ctx context.Context, userID int64) ([]string, error)
}

type historyRepository struct {
	collection *mongo.Collection
}

func NewHistoryRepository(db *mongo.Database) HistoryRepository {
	return &historyRepository{collection: db.Collection("track_listens")}
}

func (r *historyRepository) SaveTrackListen(ctx context.Context, listen model.TrackListen) error {
	_, err := r.collection.InsertOne(ctx, listen)
	return err
}

func (r *historyRepository) GetTopGenres(ctx context.Context, userID int64) ([]string, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"user_id": userID}}},
		{{Key: "$group", Value: bson.M{"_id": "$genre", "count": bson.M{"$sum": 1}}}},
		{{Key: "$sort", Value: bson.M{"count": -1}}},
		{{Key: "$limit", Value: 3}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		Genre string `bson:"_id"`
	}
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	var genres []string
	for _, r := range result {
		genres = append(genres, r.Genre)
	}
	return genres, nil
}
