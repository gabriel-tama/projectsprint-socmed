package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gabriel-tama/projectsprint-socmed/api/comment"
	"github.com/gabriel-tama/projectsprint-socmed/api/friend"
	"github.com/gabriel-tama/projectsprint-socmed/api/image"
	"github.com/gabriel-tama/projectsprint-socmed/api/post"
	"github.com/gabriel-tama/projectsprint-socmed/api/router"
	"github.com/gabriel-tama/projectsprint-socmed/api/user"
	C "github.com/gabriel-tama/projectsprint-socmed/common/config"
	psql "github.com/gabriel-tama/projectsprint-socmed/common/db"
	"github.com/gabriel-tama/projectsprint-socmed/common/jwt"
	"github.com/gin-gonic/gin"
)

func main() {

	env, err := C.Get()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, dbErr := psql.Init(context.Background())
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer db.Close(context.Background())

	// Repository/Models
	userRepository := user.NewRepository(db, env.BCRYPT_Salt)
	friendRepository := friend.NewRepository(db)
	postRepository := post.NewRepository(db)
	commentRepository := comment.NewRepository(db)

	// Services
	s3Service := image.NewS3Service(env.S3ID, env.S3Secret, env.S3Bucket, env.S3Url, env.S3Region)
	jwtService := jwt.NewJWTService(env.JWTSecret, env.JWTExp)
	userService := user.NewService(userRepository, jwtService)
	friendService := friend.NewService(friendRepository, jwtService)
	postService := post.NewService(postRepository, jwtService)
	commentService := comment.NewService(commentRepository, jwtService)

	// Controllers
	imgController := image.NewImageController(s3Service)
	userController := user.NewController(userService)
	friendControler := friend.NewController(friendService)
	postController := post.NewController(postService)
	commentController := comment.NewController(commentService)

	router := router.SetupRouter(router.RouterParam{
		JwtService:        &jwtService,
		ImageController:   imgController,
		UserController:    userController,
		FriendController:  friendControler,
		PostController:    postController,
		CommentController: commentController,
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err := router.Run(fmt.Sprintf("%s:%s", env.Host, env.Port)); err != nil {
		log.Fatal("Server error:", err)
	}
}
