package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"music_app/internal/music/model"
)

type HistoryRepository interface {
	Save(ctx context.Context, history *model.ListeningHistory) error
	GetTopGenres(ctx context.Context, userID int64) ([]string, error)
}

type historyRepo struct {
	col *mongo.Collection
}

func NewHistoryRepository(db *mongo.Database) HistoryRepository {
	return &historyRepo{col: db.Collection("listening_history")}
}

func (r *historyRepo) Save(ctx context.Context, history *model.ListeningHistory) error {
	_, err := r.col.InsertOne(ctx, history)
	return fmt.Errorf("failed to insert listening history: %w", err)
}

func (r *historyRepo) GetTopGenres(ctx context.Context, userID int64) ([]string, error) {
	cursor, err := r.col.Aggregate(ctx, mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"user_id": userID}}},
		{{Key: "$group", Value: bson.M{"_id": "$genre", "count": bson.M{"$sum": 1}}}},
		{{Key: "$sort", Value: bson.M{"count": -1}}},
		{{Key: "$limit", Value: 1}},
	})
	if err != nil {
		return nil, err
	}
	var results []struct {
		ID string `bson:"_id"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	var genres []string
	for _, g := range results {
		genres = append(genres, g.ID)
	}
	return genres, nil
}
