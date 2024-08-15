package api

import (
	"github.com/driif/golang-test-task/internal/api/handlers"
	"github.com/driif/golang-test-task/internal/server"
)

func InitRouter(s *server.Server) {
	h := handlers.NewHandlers(s)

	s.Gin.GET("/test", h.CheckHealth)
	s.Gin.POST("/message", h.PostMessage)
	s.Gin.GET("/message/list", h.GetMessages)
}
