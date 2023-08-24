package main

import (
	"log"
	"senpainikolay/go-internship-smartdata/internal/controller"
	"senpainikolay/go-internship-smartdata/internal/models"
	"senpainikolay/go-internship-smartdata/internal/repository"
	"senpainikolay/go-internship-smartdata/internal/service"
	"senpainikolay/go-internship-smartdata/pkg/db/postgres"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	router := gin.Default()
	router.Use(gin.Recovery())

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userCtrl := controller.NewUserController(userService)

	userRouter := router.Group("/user")
	{
		userRouter.GET("/me/:id", userCtrl.GetById)
		userRouter.POST("/register", userCtrl.Register)
		userRouter.POST("/login", userCtrl.LogIn)
		userRouter.DELETE("/unregister/:id", userCtrl.DeleteById)
		userRouter.PUT("/:id/img", userCtrl.UpdateImage)
		userRouter.GET("/:id/img", userCtrl.GetImage)

	}

	_ = router.Run(":8888")
}

func init() {
	db = postgres.NewDBConnection()
	err := db.AutoMigrate(models.UserModel{})
	if err != nil {
		log.Fatalf("failed to migrate user model\n")
	}
}
