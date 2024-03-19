package comment

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) CreateComment(ctx *gin.Context) {
	var req CreateCommentPayload
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	err := c.service.Create(ctx, &req)
	if errors.Is(err, ErrNotFriends) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if errors.Is(err, ErrInvalidToken) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid tkoeen"})
		return
	}

	if errors.Is(err, ErrInvalidPost) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "post not found"})
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "comment created"})

}
