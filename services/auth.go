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
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type AuthService struct {
	repository models.AuthRepository
}

func getGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
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
		"exp": time.Now().Add(24 * time.Hour).Unix(),
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
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token, err := utils.GenerateJWT(claims, jwt.SigningMethodHS256, os.Getenv("JWT_SECRET"))
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) Logout(ctx context.Context, userID uint) error {
	config := config.NewEnvConfig()
	rdb := redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
		DB:   config.RedisDB,
	})

	return rdb.Del(ctx, fmt.Sprintf("%d", userID)).Err()
}

func (s *AuthService) GenerateGoogleOAuthUrl(state string) string {
	conf := getGoogleOAuthConfig()
	return conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *AuthService) HandleGoogleCallback(ctx context.Context, code string) (string, *models.User, error) {
	conf := getGoogleOAuthConfig()

	token, err := conf.Exchange(ctx, code)
	if err != nil {
		return "", nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}

	client := conf.Client(ctx, token)
	oauth2Service, err := oauth2api.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		return "", nil, fmt.Errorf("failed to create oauth2 service: %v", err)
	}

	userinfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return "", nil, fmt.Errorf("failed to get userinfo: %v", err)
	}

	user, err := s.repository.GetUser(ctx, "email = ?", userinfo.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &models.User{
				Email:   userinfo.Email,
				Picture: userinfo.Picture,
			}
			user, err = s.repository.RegisterUserOAuth(ctx, user)
			if err != nil {
				return "", nil, fmt.Errorf("failed to register user: %v", err)
			}
		} else {
			return "", nil, fmt.Errorf("failed to query user: %v", err)
		}
	}

	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	jwtToken, err := utils.GenerateJWT(claims, jwt.SigningMethodHS256, os.Getenv("JWT_SECRET"))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create jwtToken: %v", err)
	}

	return jwtToken, user, nil
}

func NewAuthService(repository models.AuthRepository) models.AuthService {
	return &AuthService{
		repository: repository,
	}
}
