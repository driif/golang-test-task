package main

import (
	"encoding/json"
	"log"

	"github.com/driif/golang-test-task/internal/api/types"
	"github.com/driif/golang-test-task/internal/server/config"
	"github.com/go-redis/redis"
	"github.com/streadway/amqp"
)

func main() {
	cfg := config.DefaultServiceConfigFromEnv()

	// Connect to RabbitMQ
	conn, err := amqp.Dial(cfg.Rabbitmq.Addr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue
	_, err = ch.QueueDeclare(
		"Messages", // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
		DB:   cfg.Redis.DB,
	})

	// Consume messages from RabbitMQ
	msgs, err := ch.Consume(
		"Messages", // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var msg types.Msg
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			// Push the message to Redis
			err = rdb.Set(msg.Sender+"_"+msg.Receiver, msg.Message, 0).Err()
			if err != nil {
				log.Printf("Failed to push message to Redis: %v", err)
				continue
			}

			log.Printf("Message pushed to Redis: %v", msg)
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	<-forever
}
