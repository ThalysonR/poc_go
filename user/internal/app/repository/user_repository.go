package repository

import (
	"context"

	"github.com/thalysonr/poc_go/user/internal/app/model"
)

type UserRepository interface {
	Create(context context.Context, user model.User) (uint, error)
	Delete(context context.Context, id uint) error
	FindAll(context context.Context) ([]model.User, error)
	FindOne(context context.Context, id uint) (*model.User, error)
}
