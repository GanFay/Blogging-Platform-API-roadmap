package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteBlog(c *gin.Context) {
	id := c.Param("id")

	cmdTag, err := h.DB.Exec(c.Request.Context(), `DELETE FROM posts WHERE id=$1`, id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	if cmdTag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "post not found"})
	}

	c.JSON(204, gin.H{"message": "deleted post successfully"})
}
