package post

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

func (c *Controller) CreatePost(ctx *gin.Context) {
	var req CreatePostPayload

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	data, err := c.service.Create(ctx, req)

	if errors.Is(err, ErrValidationFailed) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "validation error"})
		return
	}

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "post added successfully", "data": data})

}
