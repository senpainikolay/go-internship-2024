package main

import (
	"log"
	"os"
	"senpainikolay/go-internship-smartdata/internal/chat"
	amqptransport "senpainikolay/go-internship-smartdata/internal/controller/amqp-transport"
	htttptransport "senpainikolay/go-internship-smartdata/internal/controller/http-transport"
	rpctransport "senpainikolay/go-internship-smartdata/internal/controller/rpc-transport"
	"senpainikolay/go-internship-smartdata/internal/middleware"
	"senpainikolay/go-internship-smartdata/internal/models"
	"senpainikolay/go-internship-smartdata/internal/service"
	mongodb "senpainikolay/go-internship-smartdata/pkg/db/mongo"
	"senpainikolay/go-internship-smartdata/pkg/db/postgres"
	redisdb "senpainikolay/go-internship-smartdata/pkg/db/redis"

	mongodbrepo "senpainikolay/go-internship-smartdata/internal/repository/mongodb-repo"
	postgresrepo "senpainikolay/go-internship-smartdata/internal/repository/postgres-repo"
	redisrepo "senpainikolay/go-internship-smartdata/internal/repository/redis-repo"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var db *gorm.DB
var mongo_db *mongo.Client
var redis_db *redis.Client

func main() {
	config := configInit()

	router := gin.Default()
	router.Use(gin.Recovery())

	userRepo := configRepo()

	cacheRepo := redisrepo.NewCacheRepository(redis_db)
	userService := service.NewUserService(userRepo, cacheRepo)
	userCtrl := htttptransport.NewUserController(userService)

	// Chat  Handelr
	chat_handler := chat.NewChat()
	go chat_handler.Run()

	userRouter := router.Group("/user")
	{

		userRouter.GET("/me", middleware.RequireAuth, userCtrl.GetById)
		userRouter.POST("/refreshtoken", userCtrl.RefreshToken)
		userRouter.POST("/register", userCtrl.Register)
		userRouter.POST("/login", userCtrl.LogIn)
		userRouter.DELETE("/unregister/:id", userCtrl.DeleteById)
		userRouter.PUT("/img", middleware.RequireAuth, userCtrl.UpdateImage)
		userRouter.GET("/:id/img", userCtrl.GetImage)
		userRouter.GET("/ws", middleware.RequireAuth, userCtrl.UpgradeToSocket(chat_handler))

	}

	// gRPC
	go func() {
		log.Printf("starting gRPC API server...\n")
		rpctransport.Serve(userService, ":6666")
	}()

	// RabbitMQ
	log.Printf("starting consumer RabbitMq ...\n")
	amqptransport.Serve(userService)

	_ = router.Run(config.Port)
}

func configInit() models.Config {
	yamlFile := "config/config.yaml"

	data, err := os.ReadFile(yamlFile)
	if err != nil {
		panic(err)
	}
	var config models.Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}
	return config

}
func configRepo() service.IUserRepository {
	config := configInit()
	if config.Enviroment == "dev" {
		return mongodbrepo.NewUserRepository(mongo_db)
	}

	return postgresrepo.NewUserRepository(db)
}

func init() {
	//config := configInit()
	mongo_db = mongodb.NewDBConnection()
	db = postgres.NewDBConnection()
	redis_db = redisdb.NewRedisClient()
	err := db.AutoMigrate(models.UserModel{})
	if err != nil {
		log.Fatalf("failed to migrate user model\n")
	}

}
