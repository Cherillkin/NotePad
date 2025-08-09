package models

import (
	"context"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=auth.go -destination=mocks/mock.go

type AuthCredentials struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginUserResponse struct {
	ID    uint   `json:"ID"`
	Email string `json:"Email"`
}

type RegisterUserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

type AuthRepository interface {
	RegisterUser(ctx context.Context, registerData *AuthCredentials) (*User, error)
	GetUser(ctx context.Context, query interface{}, args ...interface{}) (*User, error)
	RegisterUserOAuth(ctx context.Context, user *User) (*User, error)
}

type AuthService interface {
	Login(ctx context.Context, loginData *AuthCredentials) (string, *User, error)
	Register(ctx context.Context, registerData *AuthCredentials) (string, *User, error)
	Logout(ctx context.Context, userID uint) error
	GenerateGoogleOAuthUrl(state string) string
	HandleGoogleCallback(ctx context.Context, code string) (string, *User, error)
}

func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
