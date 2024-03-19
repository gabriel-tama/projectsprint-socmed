package comment

import (
	"github.com/gabriel-tama/projectsprint-socmed/common/jwt"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ctx *gin.Context, req *CreateCommentPayload) error
}

type commentService struct {
	repository Repository
	jwtService jwt.JWTService
}

func NewService(repository Repository, jwtService jwt.JWTService) Service {
	return &commentService{repository: repository, jwtService: jwtService}
}

func (s *commentService) Create(ctx *gin.Context, req *CreateCommentPayload) error {
	headerToken := ctx.GetHeader("Authorization")
	token, err := s.jwtService.GetPayload(headerToken)
	if err != nil {
		return ErrInvalidToken
	}

	return s.repository.Create(ctx, req, token.UserID)

}
