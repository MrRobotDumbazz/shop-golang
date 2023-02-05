package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"shop/internal/repository"
	"shop/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=auth.go -destination=mocks/mockauth.go

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
	GenerateJWT(login, password string) (accessToken string, err error)
	ParseToken(tokenString, secret string) (*TokenClaims, error)
	ValidateToken(claims *TokenClaims, isRefresh bool) (models.Seller, error)
	DeleteToken(claims *TokenClaims)
	ExpireToken(claims *TokenClaims)
	GenerateRefreshJWT(seller models.Seller) (refreshToken string, err error)
	// ParseJWT(token string) (int, error)
	// DeleteJWT(token string) error
}

// type MockCache struct{}

// func (m *MockCache) SetToken(SID int, token string) {
// 	log.Println("mock called")
// }

type Cache interface {
	SetToken(ctx context.Context, SID int, token string) error
	GetToken(ctx context.Context, ID int) (string, error)
	DeleteToken(ctx context.Context, ID int)
	ExpireToken(ctx context.Context, ID int) error
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

func comparePassword(s *models.Seller, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(s.Password), []byte(password)) != nil
}

func (s *AuthService) GenerateJWT(login, password string) (accessToken string, err error) {
	seller, err := s.repository.GetUser(login)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			log.Println("Record not found")
			return "", ErrSellerNotFound
		}
		return "", err
	}
	if comparePassword(seller, password) {
		return "", ErrSellerNotFound
	}
	var accessUID string
	if accessToken, accessUID, _, err = s.createToken(seller.ID, 600, AccessSecret); err != nil {
		return "", err
	}
	log.Printf("Token is: %s", accessToken)
	cacheJSON, err := json.Marshal(models.CachedTokens{
		AccessUID: accessUID,
	})

	ctx := contextWithTimeout()
	s.redis.SetToken(ctx, seller.ID, string(cacheJSON))
	return accessToken, nil
}

func (s *AuthService) GenerateRefreshJWT(seller models.Seller) (refreshToken string, err error) {
	var refreshUID string
	if refreshToken, refreshUID, _, err = s.createToken(seller.ID, 46000, RefreshSecret); err != nil {
		return
	}
	cacheJSON, err := json.Marshal(models.CachedTokens{
		RefreshUID: refreshUID,
	})
	ctx := contextWithTimeout()
	s.redis.SetToken(ctx, seller.ID, string(cacheJSON))
	return refreshToken, nil
}

func (s *AuthService) ValidateToken(claims *TokenClaims, isRefresh bool) (models.Seller, error) {
	ctx := contextWithTimeout()
	cacheJSON, err := s.redis.GetToken(ctx, claims.SellerId)
	if err != nil {
		return models.Seller{}, err
	}
	cachedTokens := new(models.CachedTokens)
	seller := models.Seller{}
	err = json.Unmarshal([]byte(cacheJSON), cachedTokens)
	if err != nil {
		log.Print(err)
		seller = models.Seller{
			HasToken: false,
		}
		return seller, nil
	}
	var tokenUID string
	if isRefresh {
		tokenUID = cachedTokens.RefreshUID
	} else {
		tokenUID = cachedTokens.AccessUID
	}

	if err != nil || tokenUID != claims.UID {
		seller = models.Seller{
			HasToken: false,
		}
		return seller, nil
	}
	seller, err = s.repository.GetUserInID(claims.SellerId)
	if err != nil {
		log.Print(err)
		seller = models.Seller{
			HasToken: false,
		}
		return seller, nil
	}
	seller = models.Seller{
		HasToken: true,
	}
	return seller, nil
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
			ExpiresAt: exp,
		},
		SellerId: userID,
		UID:      uuid.String(),
	})
	tokenstring, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", 0, err
	}
	log.Printf("Token is %s", tokenstring)
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
		log.Printf("Token claims %v", claims)
		return claims, nil
	}
	return nil, errors.New("Error in claims")
}

func (s *AuthService) DeleteToken(claims *TokenClaims) {
	ctx := contextWithTimeout()
	s.redis.DeleteToken(ctx, claims.SellerId)
}

func (s *AuthService) ExpireToken(claims *TokenClaims) {
	ctx := contextWithTimeout()
	s.redis.ExpireToken(ctx, claims.SellerId)
}

func contextWithTimeout() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)
	// defer cancel()
	return ctx
}
