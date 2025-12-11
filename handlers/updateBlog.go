package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateBlog(c *gin.Context) {
	idstr := c.Param("id")

	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id: " + idstr})
	}

	var newBlog Blog
	err = c.ShouldBindJSON(&newBlog)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	timeNow := time.Now()

	cmdTag, err := h.DB.Exec(c.Request.Context(), `
	UPDATE posts SET title=$1, content=$2, category=$3, tags=$4, updated_at=$6 WHERE id=$5`, newBlog.Title, newBlog.Content, newBlog.Category, newBlog.Tags, id, timeNow)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update post: " + err.Error()})
		return
	}
	if cmdTag.RowsAffected() == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successfully updated blog!"})
}
