package router

import (
	"log"
	"time"

	"github.com/gabriel-tama/projectsprint-socmed/api/image"
	"github.com/gabriel-tama/projectsprint-socmed/api/user"
	"github.com/gabriel-tama/projectsprint-socmed/common/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

var (
	limit ratelimit.Limiter
)

type RouterParam struct {
	ImageController *image.ImageController
	UserController  *user.Controller
	JwtService      *jwt.JWTService
}

func leakBucket() gin.HandlerFunc {
	prev := time.Now()
	return func(ctx *gin.Context) {
		now := limit.Take()
		log.Printf("%v", now.Sub(prev))
		prev = now
	}
}

func SetupRouter(param RouterParam) *gin.Engine {
	limit = ratelimit.New(100)
	router := gin.Default()

	router.SetTrustedProxies([]string{"::1"}) // This is for reverse proxy

	router.Use(leakBucket())
	router.Use(gin.Recovery())

	// Setup API version 1 routes
	v1 := router.Group("/v1")
	{
		user.NewRouter(v1, param.UserController, param.JwtService)
		image.NewImageRouter(v1, param.ImageController, param.JwtService)
	}

	router.GET("/rate", func(c *gin.Context) {
		c.JSON(200, "rate limiting test")
	})

	return router
}
