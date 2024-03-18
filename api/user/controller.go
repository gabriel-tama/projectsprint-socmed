package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) CreateUser(ctx *gin.Context) {
	var req CreateUserPayload

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	data, err := c.service.Create(ctx, req)

	if errors.Is(err, ErrValidationFailed) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "validation error"})
		return
	}

	if errors.Is(err, ErrUsernameAlreadyExists) {
		ctx.JSON(http.StatusConflict, gin.H{"message": "email/phone already exists"})
		return
	}

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "user succesfully created", "data": data})
}
