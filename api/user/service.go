package user

import (
	"context"
	"fmt"

	"github.com/gabriel-tama/projectsprint-socmed/common/jwt"
	"github.com/gabriel-tama/projectsprint-socmed/common/password"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ctx context.Context, req CreateUserPayload) (*UserResponse, error)
	FindByCredential(ctx context.Context, req LoginUserPayload) (*UserResponse, error)
	LinkEmail(ctx *gin.Context, req LinkEmailPayload) error
	LinkPhone(ctx *gin.Context, req LinkPhonePayload) error
	UpdateAccount(ctx *gin.Context, req UpdateAccountPayload) error
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
		Password:       hashedPassword,
	}
	if req.CredentialByEmail() {
		user.EmailCredential = req.CredentialValue
	} else {
		user.PhoneNumberCredential = req.CredentialValue
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

func (s *userService) FindByCredential(ctx context.Context, req LoginUserPayload) (*UserResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	user := &User{
		CredentialType: req.CredentialType,
	}
	if req.CredentialByEmail() {
		user.EmailCredential = req.CredentialValue
	} else {
		user.PhoneNumberCredential = req.CredentialValue
	}

	err = s.repository.FindByCredential(ctx, user)
	if err != nil {
		return nil, err
	}
	match, err := password.Matches(req.Password, user.Password)

	if err != nil {
		return nil, err
	}

	if !match {
		return nil, fmt.Errorf("%w:%w", ErrWrongPassword, err)
	}
	accessToken, err := s.jwtService.CreateToken(int(user.ID))
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		Name:        user.Name,
		Phone:       user.PhoneNumberCredential,
		Email:       user.EmailCredential,
		AccessToken: accessToken,
	}, nil

}

func (s *userService) LinkEmail(ctx *gin.Context, req LinkEmailPayload) error {

	headerToken := ctx.GetHeader("Authorization")
	token, err := s.jwtService.GetPayload(headerToken)
	if err != nil {
		return ErrInvalidToken
	}

	err = req.Validate()
	if err != nil {
		return ErrValidationFailed
	}

	return s.repository.AddEmail(ctx, req.Email, token.UserID)

}

func (s *userService) LinkPhone(ctx *gin.Context, req LinkPhonePayload) error {
	headerToken := ctx.GetHeader("Authorization")
	token, err := s.jwtService.GetPayload(headerToken)
	if err != nil {
		return ErrInvalidToken
	}

	err = req.Validate()
	if err != nil {
		return ErrValidationFailed
	}

	return s.repository.AddPhone(ctx, req.Phone, token.UserID)

}

func (s *userService) UpdateAccount(ctx *gin.Context, req UpdateAccountPayload) error {
	headerToken := ctx.GetHeader("Authorization")
	token, err := s.jwtService.GetPayload(headerToken)
	if err != nil {
		return ErrInvalidToken
	}
	err = req.Validate()
	if err != nil {
		return ErrValidationFailed
	}
	return s.repository.UpdateAccount(ctx, req.Name, req.ImageURL, token.UserID)

}
