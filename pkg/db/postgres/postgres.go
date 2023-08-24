package postgres

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDBConnection() *gorm.DB {
	dsn := fmt.Sprintf(`host=%s
						user=%s 
						password=%s 
						dbname=%s 
						port=%s 
						sslmode=disable`,
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB_NAME"),
		os.Getenv("POSTGRES_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed gorm connection: %v\n", err)
	}

	log.Printf("successfully connect to postgres db...\n")
	return db
}
