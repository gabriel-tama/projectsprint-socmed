package friend

import (
	"github.com/gabriel-tama/projectsprint-socmed/common/jwt"
	"github.com/gin-gonic/gin"
)

type Service interface {
	AddFriend(ctx *gin.Context, req AddFriendPayload) error
	DeleteFriend(ctx *gin.Context, req DeleteFriendPayload) error
	GetAllFriends(ctx *gin.Context, req GetAllFriendsPayload) (FriendListResponse, error, int)
}

type friendService struct {
	repository Repository
	jwtService jwt.JWTService
}

func NewService(repository Repository, jwtService jwt.JWTService) Service {
	return &friendService{repository: repository, jwtService: jwtService}
}

func (s *friendService) AddFriend(ctx *gin.Context, req AddFriendPayload) error {
	headerToken := ctx.GetHeader("Authorization")
	token, err := s.jwtService.GetPayload(headerToken)
	if err != nil {
		return ErrInvalidToken
	}

	if req.UserId == token.UserID {
		return ErrAlreadyFriends
	}

	return s.repository.AddFriend(ctx, token.UserID, req.UserId)
}

func (s *friendService) DeleteFriend(ctx *gin.Context, req DeleteFriendPayload) error {
	headerToken := ctx.GetHeader("Authorization")
	token, err := s.jwtService.GetPayload(headerToken)
	if err != nil {
		return err
	}

	if req.UserId == token.UserID {
		return ErrInvalidUser
	}

	return s.repository.DeleteFriend(ctx, token.UserID, req.UserId)
}

func (s *friendService) GetAllFriends(ctx *gin.Context, req GetAllFriendsPayload) (FriendListResponse, error, int) {
	headerToken := ctx.GetHeader("Authorization")
	token, err := s.jwtService.GetPayload(headerToken)
	if err != nil {
		return nil, ErrInvalidToken, 0
	}

	return s.repository.GetAllFriends(ctx, token.UserID, req)
}
