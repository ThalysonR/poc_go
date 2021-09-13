package main

import (
	"context"
	"fmt"
	"time"

	"github.com/thalysonr/poc_go/user/internal/app/service"
	"github.com/thalysonr/poc_go/user/internal/datasources"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect db")
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)

	userDs := datasources.NewUserDBDatasource(db)
	userService := service.NewUserService(userDs)
	// userService.Create(model.User{
	// 	BirthDate: "27/05/1994",
	// 	Email:     "thalysonr.castro@hotmail.com",
	// 	FirstName: "Thalyson2",
	// 	LastName:  "Castro",
	// })
	users, _ := userService.FindAll(ctx)
	fmt.Println(users)
}
