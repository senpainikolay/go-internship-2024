package mongodb

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDBConnection() *mongo.Client {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://"+os.Getenv("MONGODB_HOST")+":"+os.Getenv("MONGODB_PORT")))
	if err != nil {
		panic(err)
	}
	return client
}

func MigrateCollections() {
	dbName := "test_db"
	ctx := context.TODO()

	client := NewDBConnection()
	defer client.Disconnect(ctx)

	err := client.Database(dbName).CreateCollection(ctx, dbName)
	if err != nil {
		return
	}
	err = client.Database(dbName).CreateCollection(ctx, "users")
	if err != nil {
		fmt.Println(err)
		return
	}
}
