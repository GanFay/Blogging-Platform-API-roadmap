package handlers

import (
	"blog/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) CreateBlog(c *gin.Context) {
	var newBlog models.Blog
	err := c.ShouldBindJSON(&newBlog)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
		return
	}
	userIDstr := strconv.Itoa(userID.(int))

	_, err = h.DB.Exec(c.Request.Context(), `
		INSERT INTO posts (author_id, title, content, category, tags)
		VALUES ($1, $2, $3, $4, $5)
	`, userIDstr, newBlog.Title, newBlog.Content, newBlog.Category, newBlog.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create blog: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "post created successfully"})
}

func (h *Handler) GetAllPosts(c *gin.Context) {
	term := c.Query("term")

	var (
		rows pgx.Rows
		err  error
	)

	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	if term != "" {
		query := `
			SELECT id, author_id, title, content, category, tags, created_at, updated_at
			FROM posts
			WHERE
				title ILIKE '%' || $1 || '%'
				OR content ILIKE '%' || $1 || '%'
				OR category ILIKE '%' || $1 || '%'
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4;
		`
		rows, err = h.DB.Query(c.Request.Context(), query, term, limit, offset)
	} else {
		rows, err = h.DB.Query(c.Request.Context(), `SELECT posts.* FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2;`, limit, offset)
	}

	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.AuthorID, &p.Title, &p.Content, &p.Category, &p.Tags, &p.CreatedAt, &p.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func (h *Handler) GetPoID(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id: " + idstr})
	}
	var post models.Post
	err = h.DB.QueryRow(c.Request.Context(), `SELECT * FROM posts WHERE id=$1`, id).Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.Category, &post.Tags, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "invalid id: " + idstr})
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func (h *Handler) DeleteBlog(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
		return
	}
	var post models.Post
	id := c.Param("id")

	err := h.DB.QueryRow(c.Request.Context(), `SELECT * FROM posts WHERE id=$1`, id).Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.Category, &post.Tags, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	authorID, err := strconv.Atoi(post.AuthorID)
	if err != nil {
		return
	}
	if authorID != userID.(int) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "not permission"})
		return
	}

	cmdTag, err := h.DB.Exec(c.Request.Context(), `DELETE FROM posts WHERE id=$1`, id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	if cmdTag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "post not found"})
		return
	}

	c.JSON(204, gin.H{"message": "deleted post successfully"})
}

func (h *Handler) UpdateBlog(c *gin.Context) {
	idstr := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
		return
	}
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id: " + idstr})
	}
	var post models.Post
	err = h.DB.QueryRow(c.Request.Context(), `SELECT * FROM posts WHERE id=$1`, id).Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.Category, &post.Tags, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	atoi, err := strconv.Atoi(post.AuthorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid author id: " + post.AuthorID})
		return
	}

	if atoi != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "not permission"})
		return
	}

	var newBlog models.Blog
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
