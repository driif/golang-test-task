package server

import (
	"github.com/driif/golang-test-task/internal/server/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/streadway/amqp"
)

type Server struct {
	Config   config.Server
	Gin      *gin.Engine
	Rabbitmq *Rabbitmq
	Redis    *redis.Client
}

type Rabbitmq struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

// NewServer creates a new server
func NewServer(cfg config.Server) *Server {
	return &Server{
		Config:   cfg,
		Gin:      gin.Default(),
		Rabbitmq: nil,
		Redis:    nil,
	}
}

// InitRabbitmq initializes RabbitMQ connection
func (s *Server) InitRabbitmq() error {
	conn, err := amqp.Dial(s.Config.Rabbitmq.Addr)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	s.Rabbitmq = &Rabbitmq{
		Conn: conn,
		Ch:   ch,
	}

	_, err = s.Rabbitmq.Ch.QueueDeclare(
		"Messages", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)

	return err
}

// InitRedis initializes Redis connection
func (s *Server) InitRedis() error {
	s.Redis = redis.NewClient(&redis.Options{
		Addr: s.Config.Redis.Addr,
		DB:   s.Config.Redis.DB,
	})

	return nil
}

// Start starts the server
func (s *Server) Start() error {
	return s.Gin.Run(s.Config.Gin.ListenAddress)
}

// Close closes the server
func (s *Server) Close() {
	if s.Rabbitmq != nil {
		s.Rabbitmq.Ch.Close()
		s.Rabbitmq.Conn.Close()
	}

}
