package entity

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentInfo struct {
	StudentName  string `json:"name" bson:"name"`
	StudentEmail string `json:"email" bson:"email"`
	StudentID    string `json:"user_id" bson:"user_id"`
}

type Class struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	ClassName string             `json:"class_name" bson:"class_name"`

	AuthorMail string `json:"author_mail" bson:"author_mail"`
	UserID     string `json:"user_id" bson:"user_id"`

	TestID []primitive.ObjectID `json:"test_id" bson:"test_id"`

	StudentsAccept []StudentInfo `json:"students_accept" bson:"students_accept"`
	StudentsWait   []StudentInfo `json:"students_wait" bson:"students_wait"`

	IsPublic  bool `json:"is_public" bson:"is_public"`
	CodeClass string

	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`

	Tags []string `json:"tags" bson:"tags"`
}

func (c *Class) Validate() error {
	if c.ID == primitive.NilObjectID {
		return errors.New("invalid Class ID")
	}

	for _, testID := range c.TestID {
		if testID == primitive.NilObjectID {
			return errors.New("invalid Test ID")
		}
	}

	return nil
}

func CreateNewClass(newClass Class) (Class, error) {
	if err := newClass.Validate(); err != nil {
		return Class{}, err
	}

	newClass.UpdatedAt = time.Now()
	return newClass, nil
}
