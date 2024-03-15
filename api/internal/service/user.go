package service

import (
	"KinotekaAPI/internal/domain"
	"KinotekaAPI/internal/storage"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type userService struct {
	s storage.UserStorage
}

const (
	salt       = "aasdf78nbvcll;l8qwo2"
	signingKey = "AKJSD67asdjb&*#@maslkd"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int64 `json:"user_id"`
}

func NewUserService(s storage.UserStorage) User {
	return &userService{
		s: s,
	}
}

func (u *userService) CreateUser(user domain.User, role string) error {
	user.Password = generatePasswordHash(user.Password)
	return u.s.CreateUser(user, role)
}

func (u *userService) GenerateToken(login, password string) (string, error) {
	user, err := u.s.GetUser(login, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	return token.SignedString([]byte(signingKey))
}

func ParseToken(accessToken string) (int64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (u *userService) IsAdmin(id int64) (bool, error) {
	roles, err := u.s.GetRole(id)
	if err != nil {
		return false, err
	}

	for _, el := range roles {
		if el.Name == "admin" {
			return true, nil
		}
	}
	return false, nil
}
