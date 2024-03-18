package friend

import (
	"github.com/gabriel-tama/projectsprint-socmed/common/jwt"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.RouterGroup, controller *Controller, jwtService *jwt.JWTService) {
	router := r.Group("/friend")

	{
		router.POST("/", controller.AddFriend)
		router.DELETE("/", controller.DeleteFriend)
	}
}
