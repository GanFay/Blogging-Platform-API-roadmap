package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Me godoc
// @Summary Get current user
// @Description Returns information about the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} info "User information"
// @Failure 400 {object} map[string]string "Database error"
// @Router /users/me [get]
func (h *Handler) Me(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetInt("user_id")

	req, err := h.Users.GetByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"id": req.ID, "username": req.Username, "email": req.Email})

}
