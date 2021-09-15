package service

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	terrors "github.com/thalysonr/poc_go/common/errors"
	"github.com/thalysonr/poc_go/common/log"
	"github.com/thalysonr/poc_go/user/internal/app/model"
	"github.com/thalysonr/poc_go/user/internal/app/repository"
)

type UserService struct {
	userRepository repository.UserRepository
	validate       *validator.Validate
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
		validate:       validator.New(),
	}
}

func (u *UserService) Create(ctx context.Context, user model.User) (uint, error) {
	err := u.validate.Struct(user)
	if err != nil {
		var fields []terrors.ValidationErrorField
		for _, err := range err.(validator.ValidationErrors) {
			fields = append(fields, terrors.ValidationErrorField{
				Error: err.Error(),
				Field: err.Field(),
			})
		}
		return 0, terrors.NewValidationError(fields)
	}

	id, err := u.userRepository.Create(ctx, user)
	if err != nil {
		log.GetLogger().Warn(ctx, "could not create user: ", err)
		return 0, terrors.NewInternalServerError()
	}
	return id, nil
}

func (u *UserService) Delete(ctx context.Context, id uint) error {
	err := u.userRepository.Delete(ctx, id)
	if err != nil {
		log.GetLogger().Warn(ctx, "could not delete user: ", err)
		return terrors.NewInternalServerError()
	}
	return nil
}

func (u *UserService) FindAll(ctx context.Context) ([]model.User, error) {
	users, err := u.userRepository.FindAll(ctx)
	if err != nil {
		log.GetLogger().Warn(ctx, "could not find all users: ", err)
		return nil, terrors.NewInternalServerError()
	}
	return users, nil
}

func (u *UserService) FindOne(ctx context.Context, id uint) (*model.User, error) {
	user, err := u.userRepository.FindOne(ctx, id)
	if err != nil {
		if !errors.Is(err, &terrors.ErrNotFound{}) {
			log.GetLogger().Warn(ctx, "could not find one user: ", err)
			return nil, terrors.NewInternalServerError()
		}
		return nil, err
	}
	return user, nil
}
