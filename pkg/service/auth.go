package service

import (
	"crypto/sha1"
	"druna_server/pkg/model"
	"druna_server/pkg/repository"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	signingKey      = "jgfdi4trgdffdgdf"
	accessTokenTTL  = 12 * time.Hour
	refreshTokenTTL = 7 * 24 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user model.User) (int, error) {
	user.PasswordHash = generatePasswordHash(user.PasswordHash)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(tokenTTL time.Duration, user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
		user.Username,
	})
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) GenerateAccessRefreshToken(username, passwordHash string) (string, string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(passwordHash))
	if err != nil {
		return "", "", err
	}

	accessToken, err := s.GenerateToken(accessTokenTTL, user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.GenerateToken(refreshTokenTTL, user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

/*
	func (s *AuthService) GenerateRefreshToken(username, passwordHash string) (string, error) {
		user, err := s.repo.GetUser(username, generatePasswordHash(passwordHash))
		if err != nil {
			return "", err
		}

		refreshToken, err := s.GenerateToken(refreshTokenTTL, user)
		if err != nil {
			return "", err
		}

		return refreshToken, nil
	}
*/
func (s *AuthService) RenewToken(username string, userid int) (string, string, error) {
	user := model.User{
		ID:       userid,
		Username: username,
	}

	newAccessToken, err := s.GenerateToken(accessTokenTTL, user)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.GenerateToken(refreshTokenTTL, user)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *AuthService) ParseToken(accessToken string) (int, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, claims.Username, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprint(hash)
}
