package handlers

import "github.com/gin-gonic/gin"

func (h *Handlers) CheckHealth(c *gin.Context) {
	c.JSON(200, "worked")
}
