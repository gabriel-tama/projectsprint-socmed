package friend

import (
	"strconv"

	"github.com/gabriel-tama/projectsprint-socmed/common/jwt"
	"github.com/gin-gonic/gin"
)

type Service interface {
	AddFriend(ctx *gin.Context, req AddFriendPayload) error
	DeleteFriend(ctx *gin.Context, req DeleteFriendPayload) error
	GetAllFriends(ctx *gin.Context, req GetAllFriendsPayload) (*FriendListResponse, int, error)
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

	userId, err := strconv.Atoi(req.UserId)
	if err != nil {
		return ErrInvalidUser
	}

	if userId == token.UserID {
		return ErrAlreadyFriends
	}

	return s.repository.AddFriend(ctx, token.UserID, userId)

}

func (s *friendService) DeleteFriend(ctx *gin.Context, req DeleteFriendPayload) error {
	headerToken := ctx.GetHeader("Authorization")
	token, err := s.jwtService.GetPayload(headerToken)
	if err != nil {
		return err
	}

	if req.UserInt == token.UserID {
		return ErrInvalidUser
	}

	return s.repository.DeleteFriend(ctx, token.UserID, req.UserInt)
}

func (s *friendService) GetAllFriends(ctx *gin.Context, req GetAllFriendsPayload) (*FriendListResponse, int, error) {
	headerToken := ctx.GetHeader("Authorization")
	token, err := s.jwtService.GetPayload(headerToken)
	if err != nil {
		return nil, 0, ErrInvalidToken
	}

	return s.repository.GetAllFriends(ctx, token.UserID, req)
}
