package handlers

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/driif/golang-test-task/internal/api/types"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetMessages(c *gin.Context) {

	sender := c.Query("sender")
	receiver := c.Query("receiver")

	if sender == "" || receiver == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sender and receiver parameters are required"})
		return
	}

	// Fetch all keys matching the sender and receiver pattern
	pattern := sender + "_" + receiver + "_*"
	keys, err := h.Server.Redis.Keys(pattern).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch keys from Redis"})
		return
	}

	var messages []types.Msg
	for _, key := range keys {
		val, err := h.Server.Redis.Get(key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch message from Redis"})
			return
		}

		var msg types.Msg
		err = json.Unmarshal([]byte(val), &msg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal message"})
			return
		}

		messages = append(messages, msg)
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Timestamp.After(messages[j].Timestamp)
	})

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
