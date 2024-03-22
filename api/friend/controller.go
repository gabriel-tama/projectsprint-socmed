package friend

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	user_id, err := strconv.Atoi(req.UserId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, res)
		return
	}
	req.UserInt = user_id
	err = c.service.AddFriend(ctx, req)

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
	ctx.JSON(http.StatusOK, res)
}

func (c *Controller) DeleteFriend(ctx *gin.Context) {
	var req DeleteFriendPayload
	var res response.ResponseBody

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Message = "bad request"
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	user_id, err := strconv.Atoi(req.UserId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, res)
		return
	}
	req.UserInt = user_id
	err = c.service.DeleteFriend(ctx, req)
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

func (c *Controller) GetAllFriendsFriend(ctx *gin.Context) {
	var req GetAllFriendsPayload
	var pagination response.Pagination
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	paramPairs := ctx.Request.URL.Query()
	for _, values := range paramPairs {
		if values[0] == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			return
		}
	}

	data, err, total := c.service.GetAllFriends(ctx, req)
	if errors.Is(err, ErrValidationFailed) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error serper"})
		return
	}
	pagination.Total = total
	pagination.Limit = req.Limit
	pagination.Offset = req.Offset

	ctx.JSON(http.StatusOK, gin.H{"message": "", "data": data, "meta": pagination})
}
