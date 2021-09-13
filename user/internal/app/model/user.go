package model

import "time"

type User struct {
	ID        uint `gorm:"primarykey"`
	BirthDate string
	Email     string `validator:"required"`
	FirstName string `validator:"required"`
	LastName  string `validator:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
