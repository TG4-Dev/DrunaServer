package service

import (
	"database/sql"
	"druna_server/pkg/model"
	"druna_server/pkg/repository"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenTTL  = 12 * time.Hour
	refreshTokenTTL = 7 * 24 * time.Hour
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

type telegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type AuthService struct {
	repo       repository.Authorization
	botToken   string
	signingKey []byte
}

func NewAuthService(repo repository.Authorization) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logrus.Fatal("JWT_SECRET environment variable is required")
	}

	return &AuthService{
		repo:       repo,
		botToken:   os.Getenv("BOT_TOKEN"),
		signingKey: []byte(secret),
	}
}

func (s *AuthService) CreateUser(user model.User) (int, error) {
	hash, err := hashPassword(user.PasswordHash)
	if err != nil {
		return 0, err
	}
	user.PasswordHash = hash
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(tokenTTL time.Duration, user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:   user.ID,
		Username: user.Username,
	})
	return token.SignedString(s.signingKey)
}

func (s *AuthService) GenerateAccessRefreshToken(username, password string) (string, string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	return s.generateTokensForUser(user)
}

func (s *AuthService) RenewToken(username string, userid int) (string, string, error) {
	user := model.User{
		ID:       userid,
		Username: username,
	}

	return s.generateTokensForUser(user)
}

func (s *AuthService) ParseToken(accessToken string) (int, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return s.signingKey, nil
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

func (s *AuthService) TelegramLogin(telegramID int64, name, username string) (string, string, error) {
	user, err := s.repo.GetUserByTelegramID(telegramID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "", "", err
		}

		telegramUsername := fmt.Sprintf("tg_%d", telegramID)
		if username != "" {
			telegramUsername = username
		}

		randomPassword, err := hashPassword(uuid.New().String())
		if err != nil {
			return "", "", err
		}

		tgID := telegramID
		newUser := model.User{
			Name:         name,
			Username:     telegramUsername,
			PasswordHash: randomPassword,
			Email:        fmt.Sprintf("%s@telegram.local", telegramUsername),
			TelegramID:   &tgID,
		}

		id, err := s.repo.CreateUser(newUser)
		if err != nil {
			return "", "", err
		}

		user = newUser
		user.ID = id
	}

	return s.generateTokensForUser(user)
}

func (s *AuthService) LoginWithTelegramInitData(initData string) (string, string, error) {
	if s.botToken == "" {
		return "", "", errors.New("BOT_TOKEN is not configured")
	}

	data, err := parseInitData(initData, s.botToken)
	if err != nil {
		return "", "", err
	}

	userJSON, ok := data["user"]
	if !ok || userJSON == "" {
		return "", "", errors.New("missing user in init data")
	}

	var tgUser telegramUser
	if err := json.Unmarshal([]byte(userJSON), &tgUser); err != nil {
		return "", "", fmt.Errorf("invalid user data: %w", err)
	}

	name := tgUser.FirstName
	if tgUser.LastName != "" {
		name = name + " " + tgUser.LastName
	}

	return s.TelegramLogin(tgUser.ID, name, tgUser.Username)
}

func (s *AuthService) generateTokensForUser(user model.User) (string, string, error) {
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

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
