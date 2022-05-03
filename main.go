package main

import (
	// "fmt"
	"github.com/gin-gonic/gin"
	"bwastartup/user"
	"bwastartup/handler"
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
  )
  
  func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	api:= router.Group("api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)

	router.Run()
	// input dari user
	// handler : mapping input dari User jadi struct input
	// service : melakukan mapping dari struct input ke struct User
	// repository
	// db

  }
