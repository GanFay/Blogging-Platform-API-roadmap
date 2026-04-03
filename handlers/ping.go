package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping godoc
// @Summary Ping server
// @Description Check if server is running
// @Tags utility
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
