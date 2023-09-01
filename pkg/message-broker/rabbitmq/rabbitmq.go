package rabbitmq

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	ampq "github.com/streadway/amqp"
)

func NewRabbitMQConn() *ampq.Connection {

	dialStr := "amqp://guest:guest@" + os.Getenv("RABBITMQ_HOST") + ":" + os.Getenv("RABBITMQ_PORT")
	conn, err := ampq.Dial(dialStr)

	if err != nil {
		log.Fatalf(err.Error())
	}

	return conn
}
