package post

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gabriel-tama/projectsprint-socmed/common/response"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	_, err := c.service.Create(ctx, req)

	if errors.Is(err, ErrValidationFailed) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "validation error"})
		return
	}

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "post added successfully"})

}

func (c *Controller) GetAllPosts(ctx *gin.Context) {
	var req GetAllPostsPayload
	var pagination response.Pagination
	if err := ctx.ShouldBind(&req); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	paramPairs := ctx.Request.URL.Query()
	for _, values := range paramPairs {
		if values[0] == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			return
		}
	}

	data, total_count, err := c.service.GetAllPosts(ctx, req)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	pagination.Limit = req.Limit
	pagination.Offset = req.Offset
	pagination.Total = total_count
	ctx.JSON(http.StatusOK, gin.H{"message": "ok", "data": data, "meta": pagination})
}
