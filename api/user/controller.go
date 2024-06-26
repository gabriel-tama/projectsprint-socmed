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

	if errors.Is(err, ErrWrongPassword) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "validation error"})
		return
	}

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "user succesfully created", "data": data})
}

func (c *Controller) LoginUser(ctx *gin.Context) {
	var req LoginUserPayload

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	data, err := c.service.FindByCredential(ctx, req)

	if errors.Is(err, ErrValidationFailed) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "validation error"})
		return
	}

	if errors.Is(err, ErrUserNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	if errors.Is(err, ErrWrongPassword) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong pass"})
		return

	}

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "server eerr"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success logged in", "data": data})
}

func (c *Controller) LinkEmail(ctx *gin.Context) {
	var req LinkEmailPayload

	headerToken := ctx.GetHeader("Authorization")
	if headerToken == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	err := c.service.LinkEmail(ctx, req)

	if errors.Is(err, ErrInvalidToken) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	if errors.Is(err, ErrEmailAlreadyExists) {
		ctx.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		return
	}

	if errors.Is(err, ErrWrongRoute) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if errors.Is(err, ErrValidationFailed) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "succesfully link email"})

}

func (c *Controller) LinkPhone(ctx *gin.Context) {
	var req LinkPhonePayload
	headerToken := ctx.GetHeader("Authorization")
	if headerToken == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err := c.service.LinkPhone(ctx, req)

	if errors.Is(err, ErrInvalidToken) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "token is missing or invalid"})
		return
	}

	if errors.Is(err, ErrPhoneAlreadyExists) {
		ctx.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		return
	}

	if errors.Is(err, ErrWrongRoute) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "succesfully link email"})

}

func (c *Controller) UpdateAccount(ctx *gin.Context) {
	var req UpdateAccountPayload
	headerToken := ctx.GetHeader("Authorization")
	if headerToken == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	err := c.service.UpdateAccount(ctx, req)

	if errors.Is(err, ErrValidationFailed) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if errors.Is(err, ErrInvalidToken) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully update user profile"})
}
