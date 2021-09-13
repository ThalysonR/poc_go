package datasources

import (
	"context"

	"github.com/thalysonr/poc_go/user/internal/app/model"
	"gorm.io/gorm"
)

type UserDBDatasource struct {
	db *gorm.DB
}

func NewUserDBDatasource(db *gorm.DB) *UserDBDatasource {
	db.AutoMigrate(&model.User{})
	return &UserDBDatasource{
		db: db,
	}
}

func (u *UserDBDatasource) Create(context context.Context, user model.User) (uint, error) {
	userEntity := &model.User{}
	*userEntity = user
	res := u.db.WithContext(context).Create(userEntity)
	if res.Error != nil {
		return 0, res.Error
	}

	return userEntity.ID, nil
}

func (u *UserDBDatasource) Delete(context context.Context, user model.User) error {
	userEntity := &model.User{}
	*userEntity = user
	res := u.db.WithContext(context).Delete(userEntity)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (u *UserDBDatasource) FindAll(context context.Context) ([]model.User, error) {
	var userEntities []model.User
	res := u.db.WithContext(context).Find(&userEntities)
	if res.Error != nil {
		return nil, res.Error
	}

	return userEntities, nil
}

func (u *UserDBDatasource) FindOne(context context.Context, id uint) (*model.User, error) {
	var userEntity model.User
	res := u.db.WithContext(context).Find(&userEntity, id)
	if res.Error != nil {
		return nil, res.Error
	}

	return &userEntity, nil
}
