package app

import (
	"errors"
	"os"

	"github.com/streadway/amqp"
)

func (a *Application) initProducer() error {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		return errors.New("RABBITMQ_URL environment variable not set")
	}

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"notifications", // имя очереди
		false,
		false,
		false,
		false,
		nil,
	)

	if a.clients == nil {
		a.clients = &Clients{}
	}
	a.clients.notificationQueue = q
	a.clients.notificationChannel = ch

	return nil
}
