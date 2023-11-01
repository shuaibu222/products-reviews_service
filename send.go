package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (app *Config) SendReviewToRabbitmq(reviewData []byte) error {
	conn, err := app.rabbitmqConnection()
	if err != nil {
		log.Println("failed to create connection with rabbitmq", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Println("failed to create channel", err)
	}

	defer conn.Close()
	defer channel.Close()

	q, err := channel.QueueDeclare(
		"reviews", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		log.Println("failed to declare a queue", err)
	}

	// Declare an exchange
	// err = channel.ExchangeDeclare(
	// 	"reviews_exchange", // Exchange name
	// 	"direct",           // Exchange type (or 'topic' if needed)
	// 	true,               // Durable
	// 	false,              // Auto-deleted
	// 	false,              // Internal
	// 	false,              // No-wait
	// 	nil,                // Arguments
	// )
	// if err != nil {
	// 	return err
	// }

	// Bind the queue to the exchange
	// err = channel.QueueBind(
	// 	q.Name,             // Queue name
	// 	"reviews",          // Routing key
	// 	"reviews_exchange", // Exchange
	// 	false,
	// 	nil,
	// )
	// if err != nil {
	// 	return err
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = channel.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        reviewData,
		})

	if err != nil {
		log.Println("failed to publish with context", err)
	}

	log.Println("Successfully published")

	return nil
}

func (app *Config) rabbitmqConnection() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq/")
	if err != nil {
		return nil, err
	}

	return conn, nil
}
