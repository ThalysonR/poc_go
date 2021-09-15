package model

import "time"

type User struct {
	ID        uint   `gorm:"primarykey"`
	BirthDate string `validate:"required"`
	Email     string `validate:"required"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
