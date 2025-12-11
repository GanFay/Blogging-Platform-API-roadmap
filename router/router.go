package router

import (
	"blog/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handlers.Handler) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", h.Ping)
	r.POST("/blog", h.CreateBlog)
	r.PUT("/blog/:id", h.UpdateBlog)
	return r
}
