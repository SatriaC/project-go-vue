package main

import (
	"bwastartup/transaction"
	"strings"
	// "fmt"
	"github.com/gin-gonic/gin"
	"bwastartup/campaign"
	"bwastartup/user"
	"bwastartup/auth"
	"bwastartup/payment"
	"bwastartup/helper"
	"bwastartup/handler"
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/dgrijalva/jwt-go"
	"net/http"
  )
  
  func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	
	campaignRepository := campaign.NewRepository(db)
	campaignService := campaign.NewService(campaignRepository)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	
	paymentService := payment.NewService()
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewService()
	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()
	router.Static("/images", "./images")
	api:= router.Group("api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email-checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	
	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	router.Run()
	// input dari user
	// handler : mapping input dari User jadi struct input
	// service : melakukan mapping dari struct input ke struct User
	// repository
	// db

  }

  func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func (c *gin.Context)  {
		authHeader := c.GetHeader("Authorization")
  
		if !strings.Contains(authHeader, "Bearer"){
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2{
		  tokenString = arrayToken[1]
		}
  
	    token, err := authService.ValidateToken(tokenString)
		if err != nil{
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid{
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetByUserID(userID)
		if err != nil{
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)

	}
  }

