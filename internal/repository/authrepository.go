package repository

import (
	"database/sql"
	"errors"
	"shop/models"
)

type Auth interface {
	CreateSeller(*models.Seller) error
	GetUser(email string) (*models.Seller, error)
	GetUserInID(id int) (models.Seller, error)
}

type AuthRepository struct {
	db *sql.DB
}

var ErrRecordNotFound = errors.New("Record not found")

func newAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) CreateSeller(s *models.Seller) error {
	if _, err := r.db.Exec("INSERT INTO shopdb.sellers (email, password) VALUES (?, ?)", s.Email, s.Password); err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) GetUser(email string) (*models.Seller, error) {
	s := &models.Seller{}
	err := r.db.QueryRow("SELECT id, password FROM shopdb.sellers WHERE email = ?",
		email).Scan(&s.ID, &s.Password)
	if err == sql.ErrNoRows {
		return nil, ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *AuthRepository) GetUserInID(id int) (models.Seller, error) {
	s := models.Seller{}
	err := r.db.QueryRow("SELECT id FROM shopdb.sellers WHERE id = ?", id).Scan(&s.ID)
	if err == sql.ErrNoRows {
		return models.Seller{}, ErrRecordNotFound
	}
	if err != nil {
		return models.Seller{}, err
	}
	return s, nil
}
