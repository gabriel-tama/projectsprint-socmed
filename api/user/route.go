package user

import (
	"github.com/gabriel-tama/projectsprint-socmed/common/jwt"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.RouterGroup, controller *Controller, jwtService *jwt.JWTService) {
	router := r.Group("/user")

	{
		router.POST("/register", controller.CreateUser)
		router.POST("/login", controller.LoginUser)
		router.POST("/link/email", controller.LinkEmail)
		router.POST("/link/phone", controller.LinkPhone)
		router.PATCH("/", controller.UpdateAccount)
	}
}
