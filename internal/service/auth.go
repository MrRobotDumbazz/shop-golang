package service

import (
	"database/sql"
	"errors"
	"regexp"
	"shop/internal/repository"
	"shop/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessSecret  = "access_secret_string"
	refreshSecret = "refresh_secret_string"
)

var (
	ErrSellerNotFound  = errors.New("Incorrect email or password")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
)

type Auth interface {
	CreateSeller(*models.Seller) error
	GenerateJWT(login, password string) (string, error)
	// ParseJWT(token string) (int, error)
	// DeleteJWT(token string) error
}

// type MockCache struct{}

// func (m *MockCache) SetToken(SID int, token string) {
// 	log.Println("mock called")
// }

type Cache interface {
	SetToken(SID int, token string)
	GetToken(ID int, token string)
}

type AuthService struct {
	repository repository.Auth
	redis      Cache
}

type tokenClaims struct {
	jwt.StandardClaims
	SellerId int    `json:"seller_id"`
	UID      string `json:"uid"`
}

func newAuthService(repository repository.Auth, redis Cache) *AuthService {
	return &AuthService{
		repository: repository,
		redis:      redis,
	}
}

func BeforeCreate(s *models.Seller) error {
	if len(s.Password) > 0 {
		enc, err := encryptString(s.Password)
		if err != nil {
			return err
		}
		s.Password = enc
	}
	return nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func validSeller(s *models.Seller) error {
	validEmail, err := regexp.MatchString(`[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`, s.Email)
	if err != nil {
		return err
	}
	if !validEmail {
		return ErrInvalidEmail
	}
	return nil
}

func (s *AuthService) CreateSeller(seller *models.Seller) error {
	err := validSeller(seller)
	if err != nil {
		return err
	}
	seller.Password, err = encryptString(seller.Password)
	if err != nil {
		return err
	}
	err = s.repository.CreateSeller(seller)
	return err
}

func (s *AuthService) GenerateJWT(login, password string) (string, error) {
	seller, err := s.repository.GetUser(login, password)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "", ErrSellerNotFound
		}
		return "", err
	}
	uuid := uuid.NewV4()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		SellerId: seller,
		UID:      uuid.String(),
	})
	tokenstring, err := token.SignedString(accessSecret)
	s.redis.SetToken(seller, tokenstring)
	return tokenstring, nil
}

func (s *AuthService) ValidateToken(claims *tokenClaims) error {
}

func (s *AuthService) createToken(userID int, expireMinutes int, secret string) (
	toke string,
	uid string,
	exp int64,
	err error,
) {
	exp = time.Now().Add(time.Minute * time.Duration(expireMinutes)).Unix()
	uuid := uuid.NewV4()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		SellerId: userID,
		UID:      uuid.String(),
	})
	tokenstring, err := token.SignedString(secret)
	return tokenstring, uuid.String(), exp, nil
}
