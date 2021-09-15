package datasources

import (
	"context"
	"errors"

	terrors "github.com/thalysonr/poc_go/common/errors"
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

func (u *UserDBDatasource) Delete(context context.Context, id uint) error {
	res := u.db.WithContext(context).Delete(&model.User{}, id)
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
	err := u.db.WithContext(context).First(&userEntity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, terrors.NewErrNotFound("user not found: %d", id)
		}
		return nil, err
	}

	return &userEntity, nil
}
