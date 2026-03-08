package handlers

import (
	"blog/auth"
	"blog/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) Login(c *gin.Context) {
	var req models.Login
	var user models.Users
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.Username) < 4 || len(req.Username) > 32 {
		c.JSON(400, gin.H{"error": "username is too short or too long"})
		return
	} else if len(req.Password) < 5 || len(req.Password) > 128 {
		c.JSON(400, gin.H{"error": "password is too short or too long"})
		return
	}
	err = h.DB.QueryRow(c.Request.Context(), `SELECT * FROM users WHERE username=$1`, req.Username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "user does not exist"})
		return
	}
	if !auth.ComparePasswords(user.PasswordHash, req.Password) {
		c.JSON(401, gin.H{"error": "wrong password"})
		return
	}

	var AccessToken string
	AccessToken, err = auth.GenerateAccessJWT(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var RefreshToken string
	RefreshToken, err = auth.GenerateRefreshJWT(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(
		"refresh_token",
		RefreshToken,
		60*60*24*7,
		"/",
		"",
		false,
		true,
	)

	c.JSON(200, gin.H{
		"access_token": AccessToken,
	})
}

func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(req.Username) < 4 || len(req.Username) > 32 {
		c.JSON(400, gin.H{"error": "username is too short or too long"})
		return
	} else if len(req.Password) < 5 || len(req.Password) > 128 {
		c.JSON(400, gin.H{"error": "password is too short or too long"})
		return
	} else if len(req.Email) < 6 || len(req.Email) > 256 {
		c.JSON(400, gin.H{"error": "email is too short or too long"})
		return
	}

	hashPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(400, gin.H{"er	ror": err.Error()})
		return
	}

	_, err = h.DB.Exec(c.Request.Context(), `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
	`, req.Username, req.Email, hashPassword)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "register successfully"})
}

func (h *Handler) Refresh(c *gin.Context) {

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(401, gin.H{"message": "no refreshToken"})
		return
	}

	userID, err := auth.ParseJWT(refreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid refresh"})
		return
	}

	access, _ := auth.GenerateAccessJWT(userID)

	c.JSON(200, gin.H{
		"access_token": access,
	})

}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
	c.JSON(200, gin.H{
		"message": "logged out",
	})
}
