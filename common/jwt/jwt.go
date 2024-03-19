package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	CreateToken(user_id int) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetPayload(tokenString string) (JWTClaims, error)
	exctractHeaderToken(headerToken string) (string, error)
}

type jwtServiceImpl struct {
	secretKey []byte
	jwtExp    int
}

type JWTClaims struct {
	UserID int `json:"user_id"`
	Exp    int `json:"exp"`
}

func NewJWTService(secretKey string, JWTExp int) JWTService {
	return &jwtServiceImpl{secretKey: []byte(secretKey), jwtExp: JWTExp}
}

var (
	ErrTokenInvalid = errors.New("invalid token")
	ErrMalformed    = errors.New("malformed token")
)

func (jwtService *jwtServiceImpl) CreateToken(user_id int) (string, error) {
	claim := JWTClaims{
		UserID: user_id,
		Exp:    int(time.Now().Add(time.Hour * 24).Unix()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claim.UserID,
		"exp":     claim.Exp,
	})

	tokenString, err := token.SignedString(jwtService.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (jwtService *jwtServiceImpl) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("%v: %v", ErrTokenInvalid, isvalid)
		}

		return []byte(jwtService.secretKey), nil
	})
}

func (jwtService *jwtServiceImpl) exctractHeaderToken(headerToken string) (string, error) {
	if headerToken == "" {
		return "", ErrTokenInvalid
	}
	const BEARER_SCHEMA = "BEARER "
	tokenString := headerToken[len(BEARER_SCHEMA):]

	return tokenString, nil
}

func (jwtService *jwtServiceImpl) GetPayload(tokenString string) (JWTClaims, error) {
	tokenString, err := jwtService.exctractHeaderToken(tokenString)
	if err != nil {
		return JWTClaims{}, ErrMalformed
	}

	token, err := jwtService.ValidateToken(tokenString)
	if err != nil {
		return JWTClaims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return JWTClaims{}, ErrTokenInvalid
	}

	ret := JWTClaims{
		UserID: int(claims["user_id"].(float64)),
		Exp:    int(claims["exp"].(float64)),
	}

	return ret, nil
}
