package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/driif/golang-test-task/internal/api"
	"github.com/driif/golang-test-task/internal/server"
	"github.com/driif/golang-test-task/internal/server/config"
	"github.com/rs/zerolog/log"
)

func main() {

	cfg := config.DefaultServiceConfigFromEnv()

	s := server.NewServer(cfg)
	api.InitRouter(s)

	if err := s.InitRabbitmq(); err != nil {
		log.Fatal().Err(err).Msg("failed to init rabbitmq")
	}

	go func() {
		if err := s.Start(); err != nil {
			log.Warn().Err(err).Msg("failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	s.Close()
}
