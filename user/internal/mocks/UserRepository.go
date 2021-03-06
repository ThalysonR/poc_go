// Code generated by mockery v2.7.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	model "github.com/thalysonr/poc_go/user/internal/app/model"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0, user
func (_m *UserRepository) Create(_a0 context.Context, user model.User) (uint, error) {
	ret := _m.Called(_a0, user)

	var r0 uint
	if rf, ok := ret.Get(0).(func(context.Context, model.User) uint); ok {
		r0 = rf(_a0, user)
	} else {
		r0 = ret.Get(0).(uint)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.User) error); ok {
		r1 = rf(_a0, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: _a0, id
func (_m *UserRepository) Delete(_a0 context.Context, id uint) error {
	ret := _m.Called(_a0, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) error); ok {
		r0 = rf(_a0, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields: _a0
func (_m *UserRepository) FindAll(_a0 context.Context) ([]model.User, error) {
	ret := _m.Called(_a0)

	var r0 []model.User
	if rf, ok := ret.Get(0).(func(context.Context) []model.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOne provides a mock function with given fields: _a0, id
func (_m *UserRepository) FindOne(_a0 context.Context, id uint) (*model.User, error) {
	ret := _m.Called(_a0, id)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(context.Context, uint) *model.User); ok {
		r0 = rf(_a0, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(_a0, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
