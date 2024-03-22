package friend

import (
	"github.com/gabriel-tama/projectsprint-socmed/common/jwt"
	"github.com/gabriel-tama/projectsprint-socmed/common/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.RouterGroup, controller *Controller, jwtService *jwt.JWTService) {
	router := r.Group("/friend")
	router.Use(middleware.AuthorizeJWT(*jwtService))

	{
		router.GET("/", controller.GetAllFriendsFriend)
		router.POST("/", controller.AddFriend)
		router.DELETE("/", controller.DeleteFriend)
	}
}
