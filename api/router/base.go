package router

import (
	"log"
	"time"

	"github.com/gabriel-tama/projectsprint-socmed/api/comment"
	"github.com/gabriel-tama/projectsprint-socmed/api/friend"
	"github.com/gabriel-tama/projectsprint-socmed/api/image"
	"github.com/gabriel-tama/projectsprint-socmed/api/post"
	"github.com/gabriel-tama/projectsprint-socmed/api/user"
	"github.com/gabriel-tama/projectsprint-socmed/common/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

var (
	limit ratelimit.Limiter
)

type RouterParam struct {
	JwtService        *jwt.JWTService
	ImageController   *image.ImageController
	UserController    *user.Controller
	FriendController  *friend.Controller
	PostController    *post.Controller
	CommentController *comment.Controller
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
	limit = ratelimit.New(1000)
	router := gin.Default()

	router.SetTrustedProxies([]string{"::1"}) // This is for reverse proxy

	router.Use(leakBucket())
	router.Use(gin.Recovery())

	// Setup API version 1 routes
	v1 := router.Group("/v1")
	{
		user.NewRouter(v1, param.UserController, param.JwtService)
		image.NewImageRouter(v1, param.ImageController, param.JwtService)
		friend.NewRouter(v1, param.FriendController, param.JwtService)
		post.NewRouter(v1, param.PostController, param.JwtService)
		comment.NewRouter(v1, param.CommentController, param.JwtService)
	}

	router.GET("/rate", func(c *gin.Context) {
		c.JSON(200, "rate limiting test")
	})

	return router
}
