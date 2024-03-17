package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	C "github.com/gabriel-tama/projectsprint-socmed/common/config"
	psql "github.com/gabriel-tama/projectsprint-socmed/common/db"
	"github.com/gin-gonic/gin"
)

func main() {

	env, err := C.Get()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := env.JWTSecret
	fmt.Println(secretKey)
	dbErr := psql.Init(context.Background())
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer psql.Close(context.Background())
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
