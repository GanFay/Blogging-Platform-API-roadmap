package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Blog struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

func (h *Handler) CreateBlog(c *gin.Context) {
	var newBlog Blog
	err := c.ShouldBindJSON(&newBlog)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.DB.Exec(c.Request.Context(), `
		INSERT INTO posts (title, content, category, tags)
		VALUES ($1, $2, $3, $4)
	`, newBlog.Title, newBlog.Content, newBlog.Category, newBlog.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Blog created successfully"})
}
