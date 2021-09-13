package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/thalysonr/poc_go/common/errors"
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

func (u *UserService) Create(context context.Context, user model.User) (uint, error) {
	err := u.validate.Struct(user)
	if err != nil {
		var fields []errors.ValidationErrorField
		for _, err := range err.(validator.ValidationErrors) {
			fields = append(fields, errors.ValidationErrorField{
				Error: err.Error(),
				Field: err.Field(),
			})
		}
		return 0, errors.NewValidationError(fields)
	}

	return u.userRepository.Create(context, user)
}

func (u *UserService) Delete(context context.Context, user model.User) error {
	return u.userRepository.Delete(context, user)
}

func (u *UserService) FindAll(context context.Context) ([]model.User, error) {
	log.GetLogger().Info("FindAll called...")
	return u.userRepository.FindAll(context)
}

func (u *UserService) FindOne(context context.Context, id uint) (*model.User, error) {
	return u.userRepository.FindOne(context, id)
}
