package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gabriel-tama/projectsprint-socmed/api/image"
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

	dbErr := psql.Init(context.Background())
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer psql.Close(context.Background())

	// Repository/Models

	userRepository := user.NewRepository(psql.PgPool, env.BCRYPT_Salt)

	// Services
	s3Service := image.NewS3Service(env.S3ID, env.S3Secret, env.S3Bucket, env.S3Url)
	jwtService := jwt.NewJWTService(env.JWTSecret, env.JWTExp)
	userService := user.NewService(userRepository, jwtService)

	// Controllers
	imgController := image.NewImageController(s3Service)
	userController := user.NewController(userService)

	router := router.SetupRouter(router.RouterParam{
		JwtService:      &jwtService,
		ImageController: imgController,
		UserController:  userController,
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
