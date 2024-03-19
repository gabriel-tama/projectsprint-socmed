package post

import (
	"fmt"

	"github.com/gabriel-tama/projectsprint-socmed/common/jwt"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ctx *gin.Context, req CreatePostPayload) (*PostResponse, error)
}

type postService struct {
	repository Repository
	jwtService jwt.JWTService
}

func NewService(repository Repository, jwtService jwt.JWTService) Service {
	return &postService{repository: repository, jwtService: jwtService}
}

func (s *postService) Create(ctx *gin.Context, req CreatePostPayload) (*PostResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	headerToken := ctx.GetHeader("Authorization")
	token, err := s.jwtService.GetPayload(headerToken)
	if err != nil {
		return nil, err
	}

	post := &Post{
		UserId:     uint64(token.UserID),
		PostInHtml: req.PostInHTML,
		Tags:       req.Tags,
	}
	err = s.repository.Create(ctx, post)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &PostResponse{
		PostInHtml: req.PostInHTML,
		Tags:       req.Tags,
	}, nil

}
