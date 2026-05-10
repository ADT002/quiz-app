package repository

import (
	"context"
	entity "quiz-app/internal/domain/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	CreateUser(
		ctx context.Context,
		user entity.User) (primitive.ObjectID, error)

	GetAllUser(
		ctx context.Context) ([]entity.User, error)

	UpdateUser(
		ctx context.Context,
		user entity.User) (entity.User, error)

	Login(
		ctx context.Context,
		user entity.User) (entity.User, error)

	DeleteUser(
		ctx context.Context,
		id primitive.ObjectID) error

	GetUser(
		ctx context.Context,
		user_id string) (entity.User, error)

	UpdateByEmailID(
		ctx context.Context,
		emailID string,
		update bson.M,
	) (entity.User, error)
}
