package main

import (
	"log"
	"senpainikolay/go-internship-smartdata/internal/controller"
	"senpainikolay/go-internship-smartdata/internal/middleware"
	"senpainikolay/go-internship-smartdata/internal/models"
	"senpainikolay/go-internship-smartdata/internal/repository"
	"senpainikolay/go-internship-smartdata/internal/service"
	"senpainikolay/go-internship-smartdata/pkg/db/postgres"
	redisdb "senpainikolay/go-internship-smartdata/pkg/db/redis"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var db *gorm.DB
var redis_db *redis.Client

func main() {
	router := gin.Default()
	router.Use(gin.Recovery())

	userRepo := repository.NewUserRepository(db)
	cacheRepo := repository.NewCacheRepository(redis_db)
	userService := service.NewUserService(userRepo, cacheRepo)
	userCtrl := controller.NewUserController(userService)

	userRouter := router.Group("/user")
	{
		userRouter.GET("/me", middleware.RequireAuth, userCtrl.GetById)
		userRouter.POST("/register", userCtrl.Register)
		userRouter.POST("/login", userCtrl.LogIn)
		userRouter.DELETE("/unregister/:id", userCtrl.DeleteById)
		userRouter.PUT("/img", middleware.RequireAuth, userCtrl.UpdateImage)
		userRouter.GET("/:id/img", userCtrl.GetImage)

	}

	_ = router.Run(":8888")
}

func init() {
	db = postgres.NewDBConnection()
	redis_db = redisdb.NewRedisClient()
	err := db.AutoMigrate(models.UserModel{})
	if err != nil {
		log.Fatalf("failed to migrate user model\n")
	}
}
