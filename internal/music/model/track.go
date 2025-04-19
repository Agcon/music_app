package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Track struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title      string             `bson:"title" json:"title"`
	Artist     string             `bson:"artist" json:"artist"`
	Genre      string             `bson:"genre,omitempty" json:"genre,omitempty"`
	UploadedAt time.Time          `bson:"uploaded_at" json:"uploaded_at"`
	FileID     primitive.ObjectID `bson:"file_id" json:"file_id"`
}
