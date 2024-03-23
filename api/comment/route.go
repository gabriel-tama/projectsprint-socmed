package comment

import (
	"github.com/gabriel-tama/projectsprint-socmed/common/jwt"
	"github.com/gabriel-tama/projectsprint-socmed/common/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.RouterGroup, controller *Controller, jwtService *jwt.JWTService) {
	router := r.Group("/post/comment")
	router.Use(middleware.AuthorizeJWT(*jwtService))

	{
		router.POST("/", controller.CreateComment)
	}
}
