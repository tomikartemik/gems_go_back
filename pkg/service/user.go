package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"gems_go_back/pkg/model"
	"gems_go_back/pkg/repository"
	"gems_go_back/pkg/schema"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"math/rand"
	"os"
	"time"
)

var salt = os.Getenv("SALT")
var signingKey = os.Getenv("SIGN_KEY_STRING")

type AuthService struct {
	repo repository.User
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}

type SignInResponse struct {
	Token string               `json:"token"`
	User  schema.UserWithItems `json:"user"`
}

func NewAuthService(repo repository.User) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user model.User) (schema.ShowUser, error) {
	user.Id = uuid.Must(uuid.NewV4()).String()
	user.Password = generatePasswordHash(user.Password)
	user.IsActive = true
	if user.IsAdmin != true {
		user.IsAdmin = false
	}
	user.BestItemId = 0
	rand.Seed(time.Now().UnixNano())
	photoInt := rand.Intn(82) + 1
	user.Photo = photoInt
	return s.repo.CreateUser(user)
}

func (s *AuthService) GetUserById(id string) (schema.UserWithItems, error) {
	var userWithInventory schema.UserWithItems

	user, err := s.repo.GetUserById(id)
	if err != nil {
		return userWithInventory, err
	}

	inventory, err := s.repo.GetUserInventory(id)
	if err != nil {
		return userWithInventory, err
	}

	for i, j := 0, len(inventory)-1; i < j; i, j = i+1, j-1 {
		inventory[i], inventory[j] = inventory[j], inventory[i]
	}
	userWithInventory.ID = user.ID
	userWithInventory.Balance = user.Balance
	userWithInventory.Username = user.Username
	userWithInventory.Email = user.Email
	userWithInventory.IsActive = user.IsActive
	userWithInventory.BestItem = user.BestItem
	userWithInventory.Photo = user.Photo
	userWithInventory.Items = inventory

	return userWithInventory, nil
}
func (s *AuthService) UpdateUser(id string, user schema.InputUser) (schema.ShowUser, error) {
	return s.repo.UpdateUser(id, user)
}

func (s *AuthService) ParseToken(accesToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accesToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	return claims.UserId, nil
}

func (s *AuthService) SignIn(email string, password string) (SignInResponse, error) {
	var signInResponse SignInResponse
	var userWithItems schema.UserWithItems
	var user schema.ShowUser

	user, err := s.repo.SignIn(email, generatePasswordHash(password))
	if err != nil {
		return signInResponse, err
	}

	userWithItems, err = s.GetUserById(user.ID)
	if err != nil {
		return signInResponse, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
	})

	signInResponse.Token, _ = token.SignedString([]byte(signingKey))
	signInResponse.User = userWithItems
	return signInResponse, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) SellItem(userId string, userItemId int) (schema.UserWithItems, error) {
	var user schema.UserWithItems
	err := s.repo.SellItem(userId, userItemId)
	if err != nil {
		return user, err
	}
	user, _ = s.GetUserById(userId)
	return user, nil
}

func (s *AuthService) SellAllItems(userId string) error {
	err := s.repo.SellAllItem(userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) ChangeAvatar(userId string) (int, error) {
	rand.Seed(time.Now().UnixNano())
	newPhoto := rand.Intn(82) + 1
	err := s.repo.ChangeAvatar(userId, newPhoto)
	if err != nil {
		return 0, err
	}
	return newPhoto, nil
}
