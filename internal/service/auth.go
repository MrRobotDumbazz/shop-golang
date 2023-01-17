package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"shop/internal/repository"
	"shop/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessSecret  = "access_secret_string"
	RefreshSecret = "refresh_secret_string"
)

var (
	ErrSellerNotFound  = errors.New("Incorrect email or password")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
)

type Auth interface {
	CreateSeller(*models.Seller) error
	GenerateJWT(login, password string) (accessToken string, refreshToken string, exp int64, err error)
	ParseToken(tokenString, secret string) (*TokenClaims, error)
	ValidateToken(claims *TokenClaims, isRefresh bool) error
	DeleteToken(claims *TokenClaims)
	ExpireToken(claims *TokenClaims)
	// ParseJWT(token string) (int, error)
	// DeleteJWT(token string) error
}

// type MockCache struct{}

// func (m *MockCache) SetToken(SID int, token string) {
// 	log.Println("mock called")
// }

type Cache interface {
	SetToken(SID int, token string)
	GetToken(ID int) (string, error)
	DeleteToken(ID int)
	ExpireToken(ID int)
}

type AuthService struct {
	repository repository.Auth
	redis      Cache
}

type TokenClaims struct {
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

func (s *AuthService) GenerateJWT(login, password string) (accessToken string, refreshToken string, exp int64, err error) {
	seller, err := s.repository.GetUser(login, password)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "", "", 0, ErrSellerNotFound
		}
		return "", "", 0, err
	}
	var accessUID, refreshUID string
	if accessToken, accessUID, exp, err = s.createToken(seller, 600, AccessSecret); err != nil {
		return
	}

	if refreshToken, refreshUID, _, err = s.createToken(seller, 46000, RefreshSecret); err != nil {
		return
	}

	cacheJSON, err := json.Marshal(models.CachedTokens{
		AccessUID:  accessUID,
		RefreshUID: refreshUID,
	})
	s.redis.SetToken(seller, string(cacheJSON))
	return accessToken, refreshToken, exp, nil
}

func (s *AuthService) ValidateToken(claims *TokenClaims, isRefresh bool) error {
	cacheJSON, err := s.redis.GetToken(claims.SellerId)
	if err != nil {
		return err
	}
	cachedTokens := new(models.CachedTokens)
	err = json.Unmarshal([]byte(cacheJSON), cachedTokens)
	if err != nil {
		return nil
	}
	var tokenUID string
	if isRefresh {
		tokenUID = cachedTokens.RefreshUID
	} else {
		tokenUID = cachedTokens.AccessUID
	}

	if err != nil || tokenUID != claims.UID {
		return errors.New("token not found")
	}

	return nil
}

func (s *AuthService) createToken(userID int, expireMinutes int, secret string) (
	toke string,
	uid string,
	exp int64,
	err error,
) {
	exp = time.Now().Add(time.Minute * time.Duration(expireMinutes)).Unix()
	uuid := uuid.NewV4()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
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

func (s *AuthService) ParseToken(tokenString, secret string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Error in claims")
}

func (s *AuthService) DeleteToken(claims *TokenClaims) {
	s.redis.DeleteToken(claims.SellerId)
}

func (s *AuthService) ExpireToken(claims *TokenClaims) {
	s.redis.ExpireToken(claims.SellerId)
}
