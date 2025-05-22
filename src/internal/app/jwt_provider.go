package app

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtProvider struct{
	accessSecret []byte
	refreshSecret []byte
}

func NewJwtProvider() *JwtProvider {
	return &JwtProvider{
		accessSecret: []byte(os.Getenv("JWT_ACCESS_SECRET")), // Replace with your actual secret],
		refreshSecret: []byte(os.Getenv("JWT_REFRESH_SECRET")),
	}
}

func (j *JwtProvider) GenerateAccessToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"uuid":  user.UUID.String(),	
		"exp": time.Now().Add(time.Minute *15).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.accessSecret)
}

func (j *JwtProvider) GenerateRefreshToken(user User) (string, error) {
		claims := jwt.MapClaims{
		"uuid":  user.UUID.String(),	
		"exp": time.Now().Add(time.Hour*24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.refreshSecret)
}

func (j *JwtProvider) ValidateAccessToken(tokenStr string) (uuid.UUID, error) {
    
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
        return j.accessSecret, nil
    })
    if err != nil || !token.Valid {
        return uuid.Nil, errors.New("invalid access token")
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return uuid.Nil, errors.New("invalid claims")
    }
    id, err := uuid.Parse(claims["uuid"].(string))
    if err != nil {
        return uuid.Nil, err
    }
    return id, nil
}

func (j *JwtProvider) ValidateRefreshToken(tokenStr string) (uuid.UUID, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
        return j.refreshSecret, nil
    })
    if err != nil || !token.Valid {
        return uuid.Nil, errors.New("invalid access token")
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return uuid.Nil, errors.New("invalid claims")
    }
    id, err := uuid.Parse(claims["uuid"].(string))
    if err != nil {
        return uuid.Nil, err
    }
    return id, nil
}