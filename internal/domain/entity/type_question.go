package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type TypeQuestion struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	AuthorID primitive.ObjectID `bson:"author_id,omitempty" json:"author_id,omitempty"`
	UserID   string             `bson:"user_id,omitempty" json:"user_id,omitempty"`
	TypeName string             `bson:"type_name" json:"type_name,omitempty"`
	TypeID   int32              `bson:"type_id" json:"type_id,omitempty"`
}
