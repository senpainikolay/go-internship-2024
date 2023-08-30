package main

import (
	"encoding/json"
	"log"
	"os"
	"senpainikolay/go-internship-smartdata/internal/models"
	"senpainikolay/go-internship-smartdata/pkg/message-broker/rabbitmq"

	"github.com/streadway/amqp"
)

func main() {
	registerClient()
}

func registerClient() {
	conn := rabbitmq.NewRabbitMQConn()
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"Register", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf(err.Error())
	}

	file, err := os.Open("./client/rabbitmq/registration.json")
	if err != nil {
		log.Fatalf(err.Error())

	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var usr models.UserCredentials
	err = decoder.Decode(&usr)
	if err != nil {
		log.Fatalf(err.Error())
	}

	jsonStr, err := json.Marshal(usr)
	if err != nil {
		log.Fatalf(err.Error())
	}

	body := string(jsonStr)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Println("Succesfuly registration sent")
}
