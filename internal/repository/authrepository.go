package repository

import (
	"database/sql"
	"shop/models"
)

type Auth interface {
	CreateSeller(*models.Seller) error
	FindByLogin(login string) (*models.Seller, error)
	CreateJWT() error
}

type AuthRepository struct {
	db *sql.DB
}

func newAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) CreateSeller(s *models.Seller) error {
	if _, err := r.db.Exec("INSERT INTO shopdb.seller (email, password) VALUES (?, ?)", s.Email, s.Password); err != nil {
		return err
	}
	return nil
}
