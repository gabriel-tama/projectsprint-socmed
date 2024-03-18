package friend

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

func (c *Controller) AddFriend(ctx *gin.Context) {
	var req AddFriendPayload
	var res response.ResponseBody

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Message = "bad request"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	err := c.service.AddFriend(ctx, req)

	if errors.Is(err, ErrInvalidUser) {
		res.Message = "invalid user"
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	if errors.Is(err, ErrInvalidToken) {
		res.Message = "invalid token"
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	if errors.Is(err, ErrValidationFailed) {
		res.Message = "validation error"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if errors.Is(err, ErrAlreadyFriends) {
		res.Message = "userId is already user's friend or adding self as friend"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err != nil {
		fmt.Println(err)
		res.Message = "server error"
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Message = "friend succesfully added"
	ctx.JSON(http.StatusCreated, res)
}

func (c *Controller) DeleteFriend(ctx *gin.Context) {
	var req DeleteFriendPayload
	var res response.ResponseBody

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Message = "bad request"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	err := c.service.DeleteFriend(ctx, req)
	if errors.Is(err, ErrValidationFailed) {
		res.Message = "validation error"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if errors.Is(err, ErrNotFriends) {
		res.Message = "not friends"
		ctx.JSON(http.StatusConflict, res)
		return
	}

	if err != nil {
		fmt.Println(err)
		res.Message = "server error"
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res.Message = "friend succesfully deleted"
	ctx.JSON(http.StatusOK, res)
}
