package router

import (
	"blog/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handlers.Handler) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", h.Ping)
	r.POST("/posts", h.CreateBlog)
	r.PUT("/posts/:id", h.UpdateBlog)
	r.DELETE("/posts/:id", h.DeleteBlog)
	return r
}
