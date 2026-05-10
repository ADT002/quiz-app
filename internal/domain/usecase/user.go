package usecase

import (
	"context"
	"errors"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUseCase struct {
	UserRepo repository.UserRepository
}

func NewUserUseCase(tr repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		UserRepo: tr,
	}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, user entity.User) (primitive.ObjectID, error) {
	return uc.UserRepo.CreateUser(ctx, user)
}

func (uc *UserUseCase) GetUser(ctx context.Context, user_id string) (entity.User, error) {
	return uc.UserRepo.GetUser(ctx, user_id)
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	return uc.UserRepo.UpdateUser(ctx, user)
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	return uc.UserRepo.DeleteUser(ctx, id)
}

func (uc *UserUseCase) GetAllUser(ctx context.Context) ([]entity.User, error) {
	return uc.UserRepo.GetAllUser(ctx)
}
func (u *UserUseCase) UpdateUserPartial(
	ctx context.Context,
	emailID string,
	req entity.UpdateUserRequest,
) (entity.User, error) {

	update := bson.M{}

	if req.FirstName != nil {
		update["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		update["last_name"] = *req.LastName
	}
	// if req.Password != nil {
	// 	hashed, err := utils.HashPassword(*req.Password)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	update["password"] = hashed
	// }

	if len(update) == 0 {
		return entity.User{}, errors.New("no fields to update")
	}

	update["updated_at"] = time.Now()

	return u.UserRepo.UpdateByEmailID(ctx, emailID, update)
}
