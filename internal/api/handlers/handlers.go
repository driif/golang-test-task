package handlers

import "github.com/driif/golang-test-task/internal/server"

type Handlers struct {
	Server *server.Server
}

func NewHandlers(server *server.Server) *Handlers {
	return &Handlers{
		Server: server,
	}
}
