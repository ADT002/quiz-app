package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Level struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID    string             `bson:"user_id,omitempty" json:"user_id,omitempty"`
	LevelName string             `bson:"level_name" json:"level_name"`
}
