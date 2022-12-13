package service

import (
	"database/sql"
	"errors"
	"regexp"
	"shop/internal/repository"
	"shop/models"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrSellerNotFound  = errors.New("Incorrect email or password")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
)

type Auth interface {
	CreateSeller(*models.Seller) error
	GenerateJWT(login, password string) (string, error)
	ParseJWT(token string) (int, error)
	DeleteJWT(token string) error
}

type AuthService struct {
	repository repository.Auth
}

func newAuthService(repository repository.Auth) *AuthService {
	return &AuthService{
		repository: repository,
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

func ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(s.Password), []byte(password)) != nil
}

func validSeller(s *models.Seller) error {
	validEmail, err := regexp.MatchString(`[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`, u.Email)
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
	token := jwt.NewWithClaims(jwt.GetSigningMethod())
}
