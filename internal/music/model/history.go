package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListeningHistory struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     int64              `bson:"user_id"`
	TrackID    primitive.ObjectID `bson:"track_id"`
	Genre      string             `bson:"genre"`
	Artist     string             `bson:"artist"`
	ListenedAt time.Time          `bson:"listened_at"`
}

type TrackListen struct {
	UserID    int64              `bson:"user_id"`
	TrackID   primitive.ObjectID `bson:"track_id"`
	Genre     string             `bson:"genre"`
	Timestamp time.Time          `bson:"timestamp"`
}
