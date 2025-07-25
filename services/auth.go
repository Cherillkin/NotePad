package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Cherillkin/Notepad/config"
	"github.com/Cherillkin/Notepad/models"
	"github.com/Cherillkin/Notepad/utils"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	repository models.AuthRepository
}

func (s *AuthService) Login(ctx context.Context, loginData *models.AuthCredentials) (string, *models.User, error) {
	config := config.NewEnvConfig()
	user, err := s.repository.GetUser(ctx, "email = ?", loginData.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, fmt.Errorf("invalid credentials")
		}
		return "", nil, err
	}

	if !models.ComparePassword(loginData.Password, user.Password) {
		return "", nil, err
	}

	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Unix(),
	}

	token, err := utils.GenerateJWT(claims, jwt.SigningMethodHS256, os.Getenv("JWT_SECRET"))

	if err != nil {
		return "", nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
		DB:   config.RedisDB,
	})

	if err := rdb.Set(ctx, strconv.Itoa(int(user.ID)), token, time.Hour*24).Err(); err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) Register(ctx context.Context, registerData *models.AuthCredentials) (string, *models.User, error) {
	if !models.ValidEmail(registerData.Email) {
		return "", nil, fmt.Errorf("Please, provide a valid email to register")
	}

	if _, err := s.repository.GetUser(ctx, "email = ?", registerData.Email); !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil, fmt.Errorf("This email already use!")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}

	registerData.Password = string(hashedPassword)

	user, err := s.repository.RegisterUser(ctx, registerData)
	if err != nil {
		return "", nil, err
	}

	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Unix(),
	}

	token, err := utils.GenerateJWT(claims, jwt.SigningMethodHS256, os.Getenv("JWT_SECRET"))
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func NewAuthService(repository models.AuthRepository) models.AuthService {
	return &AuthService{
		repository: repository,
	}
}
