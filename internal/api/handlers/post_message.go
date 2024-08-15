package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/driif/golang-test-task/internal/api/types"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func (h *Handlers) PostMessage(c *gin.Context) {
	var body types.Msg
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert body to JSON
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to marshal JSON"})
		return
	}

	// Publish to RabbitMQ
	err = h.Server.Rabbitmq.Ch.Publish(
		"",         // exchange
		"Messages", // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyJSON,
		})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to publish message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
