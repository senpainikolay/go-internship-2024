package amqptransport

import (
	"encoding/json"
	"fmt"
	"log"
	"senpainikolay/go-internship-smartdata/internal/models"
	rabbitmq "senpainikolay/go-internship-smartdata/pkg/message-broker/rabbitmq"
	"time"
)

type IUserService interface {
	Register(*models.UserModel) error
}

type RegisterUserQueueController struct {
	Queue   string
	userSvc IUserService
}

func Serve(usrSvc IUserService) {

	registerUserQueue := RegisterUserQueueController{
		Queue:   "Register",
		userSvc: usrSvc,
	}
	err := initQueue(registerUserQueue.Queue)
	if err != nil {
		log.Fatalf(err.Error())
	}
	go registerUserQueue.Consume()

	log.Printf("Consumer  RabbitMQ Queues started...\n")

}

func initQueue(name string) error {
	conn := rabbitmq.NewRabbitMQConn()

	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *RegisterUserQueueController) Consume() {
	conn := rabbitmq.NewRabbitMQConn()
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		c.Queue, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)

	if err != nil {
		log.Fatalf(err.Error())

	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var usr models.UserCredentials

			err := json.Unmarshal(d.Body, &usr)
			if err != nil {
				fmt.Println("Error decoding JSON, invalid register format!")
			}
			err = c.userSvc.Register(&models.UserModel{Email: usr.Email, Password: usr.Password})
			if err != nil {
				fmt.Println("Cant register User!!!")
			}
			time.Sleep(2 * time.Second)
		}
	}()

	<-forever
}
