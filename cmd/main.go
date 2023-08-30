package main

import (
	"log"
	amqptransport "senpainikolay/go-internship-smartdata/internal/controller/amqp-transport"
	httpAPI "senpainikolay/go-internship-smartdata/internal/controller/http-transport"
	grpcAPI "senpainikolay/go-internship-smartdata/internal/controller/rpc-transport"
	"senpainikolay/go-internship-smartdata/internal/middleware"
	"senpainikolay/go-internship-smartdata/internal/models"
	"senpainikolay/go-internship-smartdata/internal/repository"
	"senpainikolay/go-internship-smartdata/internal/service"
	"senpainikolay/go-internship-smartdata/pkg/db/postgres"
	redisdb "senpainikolay/go-internship-smartdata/pkg/db/redis"
	"sync"

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
	userCtrl := httpAPI.NewUserController(userService)

	userRouter := router.Group("/user")
	{

		userRouter.GET("/me", middleware.RequireAuth, userCtrl.GetById)
		userRouter.POST("/refreshtoken", userCtrl.RefreshToken)
		userRouter.POST("/register", userCtrl.Register)
		userRouter.POST("/login", userCtrl.LogIn)
		userRouter.DELETE("/unregister/:id", userCtrl.DeleteById)
		userRouter.PUT("/img", middleware.RequireAuth, userCtrl.UpdateImage)
		userRouter.GET("/:id/img", userCtrl.GetImage)

	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		log.Printf("starting gRPC API server...\n")
		grpcAPI.Serve(userService, ":6666")
		wg.Done()
	}()

	log.Printf("starting consumer RabbitMq ...\n")
	amqptransport.Serve(userService)

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
