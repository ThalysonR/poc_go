package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thalysonr/poc_go/common/errors"
	"github.com/thalysonr/poc_go/user/internal/app/model"
	"github.com/thalysonr/poc_go/user/internal/app/service"
	"github.com/thalysonr/poc_go/user/internal/mocks"
)

func TestCreate(t *testing.T) {
	cases := map[string]struct {
		user          model.User
		createdUserId uint
	}{
		"success": {
			user: model.User{
				BirthDate: "11/09/2011",
				Email:     "always@lurking",
				FirstName: "The",
				LastName:  "Ctulhu",
			},
		},
	}

	userRepoMock := &mocks.UserRepository{}
	for caseTitle, tc := range cases {
		t.Run(caseTitle, func(t *testing.T) {
			userRepoMock.On("Create", mock.Anything, tc.user).Return(tc.createdUserId, nil)
			userSvc := service.NewUserService(userRepoMock)
			id, err := userSvc.Create(context.Background(), tc.user)

			assert.NoError(t, err, caseTitle)
			assert.Equal(t, tc.createdUserId, id, caseTitle)
		})
	}
}

func TestCreateErrors(t *testing.T) {
	cases := map[string]struct {
		user       model.User
		repoReturn error
		expected   error
	}{
		"missing-birthdate": {
			user: model.User{
				Email:     "always@lurking",
				FirstName: "The",
				LastName:  "Ctulhu",
			},
			expected: errors.NewValidationError(nil),
		},
		"missing-email": {
			user: model.User{
				BirthDate: "11/09/2011",
				FirstName: "The",
				LastName:  "Ctulhu",
			},
			expected: errors.NewValidationError(nil),
		},
		"missing-first-name": {
			user: model.User{
				BirthDate: "11/09/2011",
				Email:     "always@lurking",
				LastName:  "Ctulhu",
			},
			expected: errors.NewValidationError(nil),
		},
		"missing-last-name": {
			user: model.User{
				BirthDate: "11/09/2011",
				Email:     "always@lurking",
				FirstName: "The",
			},
			expected: errors.NewValidationError(nil),
		},
		"internal-error": {
			user: model.User{
				BirthDate: "11/09/2011",
				Email:     "always@lurking",
				FirstName: "The",
				LastName:  "Ctulhu",
			},
			repoReturn: fmt.Errorf("some error"),
			expected:   errors.NewInternalServerError(),
		},
	}

	for caseTitle, tc := range cases {
		t.Run(caseTitle, func(t *testing.T) {
			userRepoMock := &mocks.UserRepository{}
			userRepoMock.On("Create", mock.Anything, tc.user).Return(uint(0), tc.repoReturn)
			userSvc := service.NewUserService(userRepoMock)
			_, err := userSvc.Create(context.Background(), tc.user)

			assert.ErrorAs(t, err, &tc.expected, caseTitle)
		})
	}
}

func TestDelete(t *testing.T) {
	cases := map[string]struct {
		userID     uint
		repoReturn error
		expected   error
	}{
		"success": {
			userID:     uint(32),
			repoReturn: nil,
			expected:   nil,
		},
		"internal-error": {
			userID:     uint(32),
			repoReturn: fmt.Errorf("some error"),
			expected:   errors.NewInternalServerError(),
		},
	}

	for caseTitle, tc := range cases {
		t.Run(caseTitle, func(t *testing.T) {
			repoMock := &mocks.UserRepository{}
			repoMock.On("Delete", mock.Anything, tc.userID).Return(tc.repoReturn)
			svc := service.NewUserService(repoMock)
			err := svc.Delete(context.Background(), tc.userID)

			if tc.expected != nil {
				assert.ErrorAs(t, err, &tc.expected, caseTitle)
			} else {
				assert.NoError(t, err, caseTitle)
			}
		})
	}
}
