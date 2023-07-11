package controller

import (
	"context"
	"fmt"
	"os"
	"time"

	"vngitPub/pkg/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

//SendToQueue - Send POST data to RabbitMQ
func SendToQueue(body []byte) {
	logger := utils.ConfigZap()

	address := os.Getenv("ADDRESSRB")
	username := os.Getenv("USERRB")
	password := os.Getenv("PASSRB")
	port := os.Getenv("PORTRB")
	queue := os.Getenv("QUEUE")
	var _amqp string = fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, address, port)

	//Connect RabbitMQ Instance
	conn, err := amqp.Dial(_amqp)
	if err != nil {
		logger.Errorf("Connecting to RabbitMQ...FAILED: %s", err)
	} else {
		logger.Debug("Connecting to RabbitMQ...OK")
	}
	defer conn.Close()

	//Get RabbitMQ Channel
	ch, err := conn.Channel()
	if err != nil {
		logger.Errorf("Openning a channel...FAILED: %s", err)
	} else {
		logger.Debug("Openning a channel...OK")
	}
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		queue, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	if err != nil {
		logger.Errorf("Publishing a message to RabbitMQ...FAILED: %s", err)
	} else {
		logger.Debug("Publishing a message to RabbitMQ...OK")
	}
	logger.Infof("Message: %s", body)
}