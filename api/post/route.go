package post

import (
	"github.com/gabriel-tama/projectsprint-socmed/common/jwt"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.RouterGroup, controller *Controller, jwtService *jwt.JWTService) {
	router := r.Group("/post")

	{
		router.POST("/", controller.CreatePost)
	}
}
