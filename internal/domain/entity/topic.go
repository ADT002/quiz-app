package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Topic struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID    string             `bson:"user_id,omitempty" json:"user_id,omitempty"`
	TopicName string             `bson:"topic_name" json:"topic_name"`
	TopicNO   int8               `bson:"topic_no" json:"topic_no"`
}
