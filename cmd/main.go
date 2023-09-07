package main

import (
	"log"
	"os"
	tokenservice "senpainikolay/go-internship-smartdata/auth-service/_token_service/cmd"
	rpctransport "senpainikolay/go-internship-smartdata/auth-service/internal/controller/rpc-transport"
	"senpainikolay/go-internship-smartdata/auth-service/internal/models"
	"senpainikolay/go-internship-smartdata/auth-service/internal/service"
	mongodb "senpainikolay/go-internship-smartdata/auth-service/pkg/db/mongo"
	"senpainikolay/go-internship-smartdata/auth-service/pkg/db/postgres"
	redisdb "senpainikolay/go-internship-smartdata/auth-service/pkg/db/redis"

	mongodbrepo "senpainikolay/go-internship-smartdata/auth-service/internal/repository/mongodb-repo"
	postgresrepo "senpainikolay/go-internship-smartdata/auth-service/internal/repository/postgres-repo"
	redisrepo "senpainikolay/go-internship-smartdata/auth-service/internal/repository/redis-repo"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var db *gorm.DB
var mongo_db *mongo.Client
var redis_db *redis.Client

var config models.Config

func main() {
	userRepo := configRepo()
	cacheRepo := redisrepo.NewCacheRepository(redis_db)
	userService := service.NewUserService(userRepo, cacheRepo)

	go func() {
		log.Printf("starting Token Validation server...\n")
		tokenservice.RunTokenValidationService()
	}()

	log.Printf("starting gRPC API server...\n")
	rpctransport.Serve(userService, ":6666")
}

func configInit() {
	yamlFile := "config/config.yaml"
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}

}
func configRepo() service.IUserRepository {
	if config.Enviroment == "dev" {
		return mongodbrepo.NewUserRepository(mongo_db)
	}

	return postgresrepo.NewUserRepository(db)
}

func init() {
	mongo_db = mongodb.NewDBConnection()
	mongodb.MigrateCollections()
	db = postgres.NewDBConnection()
	redis_db = redisdb.NewRedisClient()
	err := db.AutoMigrate(models.UserModel{})
	if err != nil {
		log.Fatalf("failed to migrate user model\n")
	}
	configInit()

}
