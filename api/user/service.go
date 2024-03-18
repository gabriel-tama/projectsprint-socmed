package user

import (
	"context"
	"fmt"

	"github.com/gabriel-tama/projectsprint-socmed/common/jwt"
	"github.com/gabriel-tama/projectsprint-socmed/common/password"
)

type Service interface {
	Create(ctx context.Context, req CreateUserPayload) (*UserResponse, error)
}

type userService struct {
	repository Repository
	jwtService jwt.JWTService
}

func NewService(repository Repository, jwtService jwt.JWTService) Service {
	return &userService{repository: repository, jwtService: jwtService}
}

func (s *userService) Create(ctx context.Context, req CreateUserPayload) (*UserResponse, error) {

	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	hashedPassword, err := password.Hash(req.Password, s.repository.GetSalt())
	if err != nil {
		return nil, err
	}
	user := &User{
		Name:           req.Name,
		CredentialType: req.CredentialType,
		Credential:     req.CredentialValue,
		Password:       hashedPassword,
	}
	err = s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	accessToken, err := s.jwtService.CreateToken(int(user.ID))
	if err != nil {
		return nil, err
	}

	if req.CredentialByEmail() {
		return &UserResponse{
			Name:        req.Name,
			Email:       req.CredentialValue,
			AccessToken: accessToken,
		}, nil
	}

	return &UserResponse{
		Name:        req.Name,
		Phone:       req.CredentialValue,
		AccessToken: accessToken,
	}, nil
}
